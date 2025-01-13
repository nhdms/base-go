package pg_converter

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

type ProtoConverter struct {
	setter   *ProtoFieldSetter
	fieldMap map[string]fieldInfo
}

type fieldInfo struct {
	index    int
	kind     reflect.Kind
	jsonName string
	name     string
	typ      reflect.Type
}

// NewProtoConverter creates a new ProtoConverter
func NewProtoConverter(msg interface{}, location *time.Location) *ProtoConverter {
	return &ProtoConverter{
		setter:   NewProtoFieldSetter(location),
		fieldMap: buildFieldMap(msg),
	}
}

func buildFieldMap(dst interface{}) map[string]fieldInfo {
	fieldMap := make(map[string]fieldInfo)
	typ := reflect.TypeOf(dst)

	// Handle pointer type
	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
	}

	if typ.Kind() != reflect.Struct {
		return fieldMap
	}

	for i := 0; i < typ.NumField(); i++ {
		field := typ.Field(i)

		// Skip unexported fields
		if !field.IsExported() {
			continue
		}

		// Try to get json tag first
		jsonTag := field.Tag.Get("json")
		var jsonName string
		if jsonTag != "" {
			parts := strings.Split(jsonTag, ",")
			jsonName = parts[0]
		}

		// If no json tag, try db tag
		if jsonName == "" {
			if dbTag := field.Tag.Get("db"); dbTag != "" {
				parts := strings.Split(dbTag, ",")
				jsonName = parts[0]
			}
		}

		// If still no name, use field name
		if jsonName == "" {
			jsonName = strings.ToLower(field.Name)
		}

		// Skip fields explicitly marked as "-"
		if jsonName == "-" {
			continue
		}

		fieldMap[jsonName] = fieldInfo{
			index:    i,
			kind:     field.Type.Kind(),
			jsonName: jsonName,
			name:     field.Name,
			typ:      field.Type,
		}
	}

	return fieldMap
}

// ConvertToProto converts WAL message data to a protobuf message
func (c *ProtoConverter) ConvertToStruct(columnNames []string, columnTypes []string, columnValues []interface{}, msg interface{}) error {
	msgValue := reflect.ValueOf(msg).Elem()

	for i, colName := range columnNames {
		// Try exact match first
		f, ok := c.fieldMap[colName]
		if !ok {
			// If not found, try case-insensitive match
			for mapKey, info := range c.fieldMap {
				if strings.EqualFold(mapKey, colName) {
					f = info
					ok = true
					break
				}
			}
			if !ok {
				// Skip columns that don't map to protobuf fields
				continue
			}
		}

		field := msgValue.Field(f.index)
		if err := c.setter.SetFieldValue(field, columnTypes[i], columnValues[i]); err != nil {
			return fmt.Errorf("setting field %s: %v", colName, err)
		}
	}

	return nil
}
