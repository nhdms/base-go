package utils

import (
	"github.com/spf13/cast"
	"strconv"
	"strings"
)

func DistinctInt64Slice(in []int64) []int64 {
	seen := make(map[int64]bool)
	distinct := make([]int64, 0)

	for _, num := range in {
		if !seen[num] {
			seen[num] = true
			distinct = append(distinct, num)
		}
	}

	return distinct
}

func StringSliceContains(in []string, s string) bool {
	for _, str := range in {
		if str == s {
			return true
		}
	}

	return false
}

func SliceContains[T comparable](slice []T, value T) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}

func ToInt64Slice(in []string) []int64 {
	result := make([]int64, len(in))
	for i, s := range in {
		result[i] = cast.ToInt64(s)
	}
	return result
}

func Intersect[T comparable](a, b []T) []T {
	m := make(map[T]bool)
	for _, item := range a {
		m[item] = true
	}
	var result []T
	for _, item := range b {
		if m[item] {
			result = append(result, item)
		}
	}
	return result
}

func IntSliceToString(slice []int64, sep string) string {
	strs := make([]string, len(slice))
	for i, v := range slice {
		strs[i] = strconv.Itoa(int(v))
	}
	return strings.Join(strs, sep)
}

func Includes[T comparable](slice []T, value T) bool {
	for _, v := range slice {
		if v == value {
			return true
		}
	}
	return false
}
