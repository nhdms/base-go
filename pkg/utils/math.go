package utils

import "github.com/spf13/cast"

func ToPercentFixed2[T comparable](a, b T) float64 {
	if cast.ToFloat64(b) == 0 {
		return 0
	}
	return float64(int(cast.ToFloat64(a)*10000/cast.ToFloat64(b))) / 100
}

func ToFixed2[T comparable](a, b T) float64 {
	if cast.ToFloat64(b) == 0 {
		return 0
	}
	return float64(int(cast.ToFloat64(a)*100/cast.ToFloat64(b))) / 100
}
