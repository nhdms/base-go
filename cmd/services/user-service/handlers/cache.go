package handlers

import "fmt"

func GetUserByIdCacheKey(id int64) string {
	return fmt.Sprintf("get_user_by_id_%d", id)
}
