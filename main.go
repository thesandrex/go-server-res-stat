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
	maxErrors             = 3
)

func main() {
	errCount := 0

	for {
		if errCount >= maxErrors {
			fmt.Println("Unable to fetch server statistic.")
			break
		}

		response, err := http.Get(serverURL)
		if err != nil || response.StatusCode != http.StatusOK {
			errCount++
			continue
		}

		defer response.Body.Close()
		body, err := io.ReadAll(response.Body)
		if err != nil {
			errCount++
			continue
		}

		// Парсинг данных
		stats := strings.Split(string(body), ",")
		if len(stats) != 7 {
			errCount++
			continue
		}

		// Обнуление счётчика при успешной попытке
		errCount = 0

		loadAvg, _ := strconv.ParseFloat(stats[0], 64)
		totalMem, _ := strconv.ParseFloat(stats[1], 64)
		usedMem, _ := strconv.ParseFloat(stats[2], 64)
		totalDisk, _ := strconv.ParseFloat(stats[3], 64)
		usedDisk, _ := strconv.ParseFloat(stats[4], 64)
		totalNet, _ := strconv.ParseFloat(stats[5], 64)
		usedNet, _ := strconv.ParseFloat(stats[6], 64)

		// Проверка Load Average
		if loadAvg > loadAvgThreshold {
			fmt.Printf("Load Average is too high: %.0f\n", loadAvg)
		}

		// Проверка использования памяти
		memUsage := (usedMem / totalMem) * 100
		if memUsage > memoryUsageThreshold {
			fmt.Printf("Memory usage too high: %.0f%%\n", math.Floor(memUsage))
		}

		// Проверка использования дискового пространства
		freeDisk := (totalDisk - usedDisk) / (1024 * 1024) // байты -> мегабайты
		diskUsage := (usedDisk / totalDisk) * 100
		if diskUsage > diskUsageThreshold {
			fmt.Printf("Free disk space is too low: %.0f Mb left\n", math.Floor(freeDisk))
		}

		// Проверка использования сети
		freeNet := (totalNet - usedNet) / (1000 * 1000) // байты -> мегабиты
		netUsage := (usedNet / totalNet) * 100
		if netUsage > networkUsageThreshold {
			fmt.Printf("Network bandwidth usage high: %.0f Mbit/s available\n", math.Round(freeNet))
		}

		// Пауза перед следующим запросом
		time.Sleep(5 * time.Second)
	}
}
