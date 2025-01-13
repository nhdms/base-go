package utils

import (
	"reflect"
	"strings"
)

func GetListFieldByTag(dest interface{}, tagPriorities ...string) []string {
	v := reflect.ValueOf(dest)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	// Handle slice type
	if v.Kind() == reflect.Slice {
		// Get the type of slice elements
		elemType := v.Type().Elem()
		if elemType.Kind() == reflect.Ptr {
			elemType = elemType.Elem()
		}
		if elemType.Kind() != reflect.Struct {
			return nil
		}
		// Create a new instance of the struct type to parse its fields
		v = reflect.New(elemType).Elem()
	} else if v.Kind() != reflect.Struct {
		return nil
	}

	var columns []string
	if len(tagPriorities) == 0 {
		tagPriorities = []string{"json"}
	}

	for i := 0; i < v.NumField(); i++ {
		field := v.Type().Field(i)
		col := ""

		for _, tagName := range tagPriorities {
			tagVal := strings.Split(field.Tag.Get(tagName), ",")[0]
			if len(tagVal) > 0 && tagVal != "-" {
				col = tagVal
				break
			}
		}

		if len(col) == 0 {
			continue
		}

		columns = append(columns, col)
	}
	return columns
}
