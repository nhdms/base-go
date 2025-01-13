package utils

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
