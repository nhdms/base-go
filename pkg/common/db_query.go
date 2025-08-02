package common

import (
	"fmt"

	"github.com/Masterminds/squirrel"
)

func SubQueryIn[T any](property string, values []T) squirrel.Sqlizer {
	return squirrel.Expr(fmt.Sprintf("%s IN (%s)", property, SliceToQueryIn(ToInterfaceSlice(values))))
}

func SliceToQueryIn(values []interface{}) string {
	queryIn := ""
	for i, v := range values {
		if i == 0 {
			queryIn += fmt.Sprintf("%v", v)
		} else {
			queryIn += fmt.Sprintf(",%v", v)
		}
	}
	return queryIn
}

func WrapJSONbyte(jsonByte []byte) []byte {
	if jsonByte == nil {
		return []byte("{}")
	}
	return jsonByte
}
