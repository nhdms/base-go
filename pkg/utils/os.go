package utils

import (
	"os"
	"runtime"
)

func IsTestMode() bool {
	for _, args := range os.Args {
		if args == "-test.v" {
			return true
		}
	}
	return false
}

func GetAvailableMemory() uint64 {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	// You can use memStats.Sys or another appropriate memory metric.
	// Consider leaving some memory unused for other processes.
	availableMemory := memStats.Sys - memStats.HeapInuse
	return availableMemory
}
