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
	// return "22,4904081501,2226235757,484715800471,121890934712,6505107771,6443546090"
	// return "3,4973030119,4322636569,331963281376,71548907549,1850323696,455050433"

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
	if stats.LoadAvg > 30 {
		fmt.Printf("Load Average is too high: %d\n", stats.LoadAvg)
	}

	memoryUsagePercent := float64(stats.MemoryUsed) / float64(stats.MemoryAvailable) * 100
	if memoryUsagePercent > 80 {
		fmt.Printf("Memory usage too high: %.0f%%\n", memoryUsagePercent)
	}

	freeDiskSpace := (stats.DiskAvailable - stats.DiskUsed) / (1024 * 1024) // Остаток в мегабайтах
	diskUsagePercent := float64(stats.DiskUsed) / float64(stats.DiskAvailable) * 100
	if diskUsagePercent > 90 {
		fmt.Printf("Free disk space is too low: %d Mb left\n", freeDiskSpace)
	}

	networkUsagePercent := float64(stats.NetworkLoadUsed) / float64(stats.NetworkLoadAvailable) * 100
	if networkUsagePercent > 90 {
		freeNetworkBandwidth := (stats.NetworkLoadAvailable - stats.NetworkLoadUsed) * 8 / (1024 * 1024) // Свободная полоса в мегабитах
		fmt.Printf("Network bandwidth usage high: %d Mbit/s available\n", freeNetworkBandwidth)
	}
}
