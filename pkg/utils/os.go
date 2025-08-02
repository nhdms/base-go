package utils

import (
	"os"
	"runtime"
	"strings"
)

var testArgs = map[string]bool{
	"test":    true,
	"-test.v": true,
}

func IsTestMode() bool {
	if IsLocalDevelopment() {
		return true
	}
	for _, args := range os.Args {
		if testArgs[args] {
			return true
		}

		for k := range testArgs {
			if strings.HasPrefix(args, k) {
				return true
			}
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

func IsLocalDevelopment() bool {
	if os.Getenv("LOCAL_DEVELOPMENT") == "true" {
		return true
	}

	// is run inside pod
	if len(os.Getenv("HOSTNAME")) > 0 && len(os.Getenv("SERVICE_NAME")) > 0 {
		return false
	}

	return true
}
