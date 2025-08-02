package common

import (
	"fmt"
	"net/http"
)

type SortByParam struct {
	Field string `json:"field"`
	Order string `json:"order"`
}

func ExtractSortByParam(r *http.Request, defaultSortField, defaultOrder string, allowedFields []string) (SortByParam, error) {
	sortField, order := defaultSortField, defaultOrder // if both are empty, default values will be used
	queryParams := r.URL.Query()
	if sortFieldReq := queryParams.Get("sort"); sortFieldReq != "" {
		sortField = sortFieldReq
		if !InSlice(sortField, allowedFields) {
			return SortByParam{}, fmt.Errorf("invalid sort field, supported fields: %+v", allowedFields)
		}
	}
	if orderReq := queryParams.Get("order"); orderReq != "" {
		order = orderReq
		if order != "asc" && order != "desc" {
			return SortByParam{}, fmt.Errorf("invalid order, supported values: asc, desc")
		}
	} else {
		if sortField == defaultSortField {
			order = defaultOrder
		} else {
			order = "asc" // if sort field is not default, default order is asc
		}
	}

	return SortByParam{
		Field: sortField,
		Order: order,
	}, nil
}

func BoolToOrder(b bool) string {
	if b {
		return "asc"
	}
	return "desc"
}
