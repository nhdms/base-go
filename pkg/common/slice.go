package common

func InSlice(s string, slice []string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}

func ToInterfaceSlice[T any](input []T) []interface{} {
	result := make([]interface{}, len(input))
	for i, v := range input {
		result[i] = v
	}
	return result
}
