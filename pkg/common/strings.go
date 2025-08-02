package common

import (
	"fmt"
	"strconv"
	"strings"
)

func RetrieveHour(timeRange string) (int, error) {
	parts := strings.Split(timeRange, " - ")
	if len(parts) < 1 {
		return 0, fmt.Errorf("invalid time range format")
	}

	hourStr := strings.Split(parts[0], ":")[0]
	hour, err := strconv.Atoi(hourStr)
	if err != nil {
		return 0, fmt.Errorf("invalid hour format: %v", err)
	}

	return hour, nil
}
