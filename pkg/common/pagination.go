package common

import (
	"fmt"
	"net/http"

	"github.com/spf13/cast"
)

type Pagination struct {
	Page     int    `json:"page"`
	Limit    int    `json:"limit"`
	Total    int64  `json:"total,omitempty"`
	NextPage string `json:"next_page,omitempty"`
}

type PaginationQueryParam struct {
	Page  int
	Limit int
}

func ExtractPaginationQueryParam(r *http.Request, defaultPage int, defaultLimit int) (PaginationQueryParam, error) {
	page, limit := defaultPage, defaultLimit
	queryParams := r.URL.Query()
	if pageStr := queryParams.Get("page"); pageStr != "" {
		page = cast.ToInt(pageStr)
	}
	if page <= 0 {
		return PaginationQueryParam{}, fmt.Errorf("invalid page value, page must be greater than 0")
	}
	if limitStr := queryParams.Get("limit"); limitStr != "" {
		limit = cast.ToInt(limitStr)
	}
	if limit <= 0 {
		return PaginationQueryParam{}, fmt.Errorf("invalid limit value, limit must be greater than 0")
	}
	return PaginationQueryParam{
		Page:  page,
		Limit: limit,
	}, nil
}
