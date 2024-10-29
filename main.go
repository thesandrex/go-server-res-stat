package main

import (
	"fmt"
	"io"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	serverURL             = "http://srv.msk01.gigacorp.local/_stats"
	loadAvgThreshold      = 30.0
	memoryUsageThreshold  = 80.0
	diskUsageThreshold    = 90.0
	networkUsageThreshold = 90.0
)

func main() {

	for {
		stats, err := getStatistics()
		if err != nil {
			fmt.Println("Unable to fetch server statistic.")
			break
		}

		loadAvg, _ := strconv.ParseFloat(stats[0], 64)
		totalMem, _ := strconv.ParseFloat(stats[1], 64)
		usedMem, _ := strconv.ParseFloat(stats[2], 64)
		totalDisk, _ := strconv.ParseFloat(stats[3], 64)
		usedDisk, _ := strconv.ParseFloat(stats[4], 64)
		totalNet, _ := strconv.ParseFloat(stats[5], 64)
		usedNet, _ := strconv.ParseFloat(stats[6], 64)

		if loadAvg > loadAvgThreshold {
			fmt.Printf("Load Average is too high: %.0f\n", loadAvg)
		}

		memUsage := (usedMem / totalMem) * 100
		if memUsage > memoryUsageThreshold {
			fmt.Printf("Memory usage too high: %.0f%%\n", math.Floor(memUsage))
		}

		freeDisk := (totalDisk - usedDisk) / (1024 * 1024) // байты -> мегабайты
		diskUsage := (usedDisk / totalDisk) * 100
		if diskUsage > diskUsageThreshold {
			fmt.Printf("Free disk space is too low: %.0f Mb left\n", math.Floor(freeDisk))
		}

		freeNet := (totalNet - usedNet) / (1000 * 1000) // байты -> мегабиты
		netUsage := (usedNet / totalNet) * 100
		if netUsage > networkUsageThreshold {
			fmt.Printf("Network bandwidth usage high: %.0f Mbit/s available\n", math.Round(freeNet))
		}

		time.Sleep(5 * time.Second)
	}
}

func getStatistics() ([]string, error) {
	response, err := http.Get(serverURL)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	body, _ := io.ReadAll(response.Body)

	stats := strings.Split(string(body), ",")

	return stats, nil
}
