package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type Stat struct {
	LoadAvg              uint64
	MemoryAvailable      uint64
	MemoryUsed           uint64
	DiskAvailable        uint64
	DiskUsed             uint64
	NetworkLoadAvailable uint64
	NetworkLoadUsed      uint64
}

func main() {
	url := "http://srv.msk01.gigacorp.local/_stats"
	client := &http.Client{}

	resp := monitorResources(url, client)
	// fmt.Println(resp)
	statistics := mapResponse(resp)
	// fmt.Println(statistics)
	evaluateStatistics(statistics)
}

func monitorResources(url string, client *http.Client) string {
	req, err := http.NewRequest("GET", url, nil)
	throwError("Error creating request:", err)

	// test
	// return "12,4819661917,2328437675,469781716496,97858919593,1315809605,1245267076"

	resp, err := client.Do(req)
	throwError("Unable to fetch server statistic:", err)
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	throwError("Error reading response body:", err)

	if resp.StatusCode == http.StatusOK {
		return string(body)
	} else {
		return resp.Status
	}
}

func throwError(msg string, err error) {
	if err != nil {
		fmt.Println(msg, err)
	}
}

func mapResponse(response string) Stat {
	elements := strings.Split(response, ",")
	converted := convertValues(elements)

	stats := Stat{
		LoadAvg:              converted[0],
		MemoryAvailable:      converted[1],
		MemoryUsed:           converted[2],
		DiskAvailable:        converted[3],
		DiskUsed:             converted[4],
		NetworkLoadAvailable: converted[5],
		NetworkLoadUsed:      converted[6],
	}

	return stats
}

func convertValues(parts []string) []uint64 {
	numbers := make([]uint64, 0, len(parts))

	for _, part := range parts {
		num, err := strconv.ParseUint(part, 10, 64)
		throwError("Error converting to uint64:", err)
		numbers = append(numbers, num)
	}

	return numbers
}

func evaluateStatistics(stats Stat) {
	if stats.LoadAvg > 30.0 {
		fmt.Printf("Load Average is too high: %d\n", stats.LoadAvg)
	}

	memUsage := (stats.MemoryUsed / stats.MemoryAvailable) * 100
	if memUsage > 80.0 {
		fmt.Printf("Memory usage too high: %d\n", memUsage)
	}

	freeDisk := (stats.DiskAvailable - stats.DiskUsed) / (1024 * 1024) // байты -> мегабайты
	diskUsage := (stats.DiskUsed / stats.DiskAvailable) * 100
	if diskUsage > 90.0 {
		fmt.Printf("Free disk space is too low: %d Mb left\n", freeDisk)
	}

	freeNet := (stats.NetworkLoadAvailable - stats.NetworkLoadUsed) / (1000 * 1000)
	netUsage := (stats.NetworkLoadUsed / stats.NetworkLoadAvailable) * 100
	if netUsage > 90.0 {
		fmt.Printf("Network bandwidth usage high: %d Mbit/s available\n", freeNet)
	}
}
