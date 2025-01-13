package dbtool

import (
	"fmt"
	"github.com/goccy/go-json"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"reflect"
	"strings"
	"time"
)

// ScanRow scans a row into any struct or proto message
// ScanRow scans a row into any struct or proto message
func ScanRow(row sqlx.ColScanner, dest interface{}) error {
	// Get column names
	cols, err := row.Columns()
	if err != nil {
		return err
	}

	// Create slice of interfaces to scan into
	values := make([]interface{}, len(cols))
	for i := range values {
		values[i] = new(interface{})
	}

	// Scan into interfaces
	if err := row.Scan(values...); err != nil {
		return err
	}

	// Use reflection to set values
	v := reflect.ValueOf(dest).Elem()
	t := v.Type()

	// Create field name map for case-insensitive matching
	fieldNames := make(map[string]string)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		fieldNames[strings.ToLower(field.Name)] = field.Name
	}

	for i, colName := range cols {
		// Try to find matching field
		fieldName, ok := fieldNames[strings.ToLower(strings.ReplaceAll(colName, "_", ""))]
		if !ok {
			continue
		}

		field := v.FieldByName(fieldName)
		if !field.IsValid() || !field.CanSet() {
			continue
		}

		val := reflect.ValueOf(values[i]).Elem().Interface()
		if val == nil {
			continue
		}

		// Handle special cases
		switch field.Type() {
		case reflect.TypeOf(&timestamppb.Timestamp{}):
			if t, ok := val.(time.Time); ok {
				field.Set(reflect.ValueOf(timestamppb.New(t)))
			}
		case reflect.TypeOf(&structpb.Struct{}):
			var mapValue map[string]interface{}

			switch v := val.(type) {
			case string:
				if err := json.Unmarshal([]byte(v), &mapValue); err != nil {
					return fmt.Errorf("parsing JSON string to map: %v", err)
				}
			case map[string]interface{}:
				mapValue = v
			case []byte:
				if err := json.Unmarshal(v, &mapValue); err != nil {
					return fmt.Errorf("parsing JSON bytes to map: %v", err)
				}
			default:
				return fmt.Errorf("unexpected type for Struct: %T", val)
			}

			// Convert to structpb.Struct
			pbStruct, err := structpb.NewStruct(mapValue)
			if err != nil {
				return fmt.Errorf("converting to protobuf Struct: %v", err)
			}

			// Set the field value
			field.Set(reflect.ValueOf(pbStruct))
		default:
			if field.Kind() != reflect.Slice && field.Kind() != reflect.Struct {
				// Handle regular types
				rv := reflect.ValueOf(val)
				if rv.Type().ConvertibleTo(field.Type()) {
					field.Set(rv.Convert(field.Type()))
				}
				continue
			}

			// Handle struct types (for JSONB)
			if field.Kind() == reflect.Struct {
				switch v := val.(type) {
				case []byte:
					// Try to unmarshal JSON data into the struct
					newStruct := reflect.New(field.Type()).Interface()
					if err := json.Unmarshal(v, newStruct); err == nil {
						field.Set(reflect.ValueOf(newStruct).Elem())
					}
				case string:
					// Handle case where JSON might be scanned as string
					newStruct := reflect.New(field.Type()).Interface()
					if err := json.Unmarshal([]byte(v), newStruct); err == nil {
						field.Set(reflect.ValueOf(newStruct).Elem())
					}
				case map[string]interface{}:
					// Handle case where JSON is already decoded
					jsonBytes, err := json.Marshal(v)
					if err == nil {
						newStruct := reflect.New(field.Type()).Interface()
						if err := json.Unmarshal(jsonBytes, newStruct); err == nil {
							field.Set(reflect.ValueOf(newStruct).Elem())
						}
					}
				}
				continue
			}

			// Handle slice types
			switch v := val.(type) {
			case []byte: // Handle _int8, _int4, _int2 array types that come as []byte
				// Try to convert []byte to pq.Int64Array
				var int64Arr pq.Int64Array
				if err := int64Arr.Scan(v); err == nil {
					// Create a new slice of the correct type
					newSlice := reflect.MakeSlice(field.Type(), len(int64Arr), len(int64Arr))

					// Convert each element to the target type
					for j, num := range int64Arr {
						elem := reflect.ValueOf(num).Convert(field.Type().Elem())
						newSlice.Index(j).Set(elem)
					}
					field.Set(newSlice)
					continue
				}

				// Try string array if int array fails
				var strArr pq.StringArray
				if err := strArr.Scan(v); err == nil {
					// If the target field is []string
					if field.Type().Elem().Kind() == reflect.String {
						field.Set(reflect.ValueOf([]string(strArr)))
					}
				}
			case pq.Int64Array:
				// Create a new slice of the correct type
				newSlice := reflect.MakeSlice(field.Type(), len(v), len(v))

				// Convert each element to the target type
				for j, num := range v {
					elem := reflect.ValueOf(num).Convert(field.Type().Elem())
					newSlice.Index(j).Set(elem)
				}
				field.Set(newSlice)
			case pq.StringArray:
				if field.Type().Elem().Kind() == reflect.String {
					field.Set(reflect.ValueOf([]string(v)))
				}
			default:
				// Handle other array types or fall back to normal conversion
				rv := reflect.ValueOf(val)
				if rv.Type().ConvertibleTo(field.Type()) {
					field.Set(rv.Convert(field.Type()))
				}
			}
		}
	}

	return nil
}

// ScanAll scans all rows into a slice of any type
func ScanAll(rows *sqlx.Rows, dest interface{}) error {
	// Check if dest is a pointer to a slice
	value := reflect.ValueOf(dest)
	if value.Kind() != reflect.Ptr {
		return fmt.Errorf("dest must be a pointer to a slice")
	}

	sliceValue := value.Elem()
	if sliceValue.Kind() != reflect.Slice {
		return fmt.Errorf("dest must be a pointer to a slice")
	}

	// Get the type of slice elements
	elementType := sliceValue.Type().Elem()
	isPtr := elementType.Kind() == reflect.Ptr
	if isPtr {
		elementType = elementType.Elem()
	}

	for rows.Next() {
		// Create a new element
		newElement := reflect.New(elementType)

		// Scan the row into the new element
		if err := ScanRow(rows, newElement.Interface()); err != nil {
			return err
		}

		// Append to the slice
		if isPtr {
			sliceValue.Set(reflect.Append(sliceValue, newElement))
		} else {
			sliceValue.Set(reflect.Append(sliceValue, newElement.Elem()))
		}
	}

	return rows.Err()
}
