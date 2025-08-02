package pg_converter

import (
	"fmt"
	"github.com/goccy/go-json"
	"github.com/spf13/cast"
	"google.golang.org/protobuf/types/known/structpb"
	"reflect"
	"time"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// ProtoFieldSetter handles setting protobuf field values
type ProtoFieldSetter struct {
	converter *PostgreSQLTypeConverter
}

// NewProtoFieldSetter creates a new ProtoFieldSetter
func NewProtoFieldSetter(location *time.Location) *ProtoFieldSetter {
	return &ProtoFieldSetter{
		converter: NewPostgreSQLTypeConverter(location),
	}
}

// SetFieldValue sets a protobuf field value based on PostgreSQL type
func (p *ProtoFieldSetter) SetFieldValue(field reflect.Value, pgType string, value interface{}) error {
	if value == nil {
		return nil
	}

	fieldType := field.Type()

	// Handle some google.protobuf.* specifically types
	switch fieldType {
	case reflect.TypeOf(&structpb.Struct{}):
		return p.setStructField(field, pgType, value)
	case reflect.TypeOf(&timestamppb.Timestamp{}):
		return p.setTimestampField(field, pgType, value)
	}

	// Convert the value using the PostgreSQL converter
	convertedValue, err := p.converter.ConvertValue(pgType, value)
	if err != nil {
		return fmt.Errorf("converting value: %v", err)
	}

	return p.setField(field, fieldType, convertedValue)
}

func (p *ProtoFieldSetter) setField(field reflect.Value, fieldType reflect.Type, value interface{}) error {
	switch fieldType.Kind() {
	case reflect.Int32, reflect.Int64, reflect.Int8, reflect.Int16:
		field.SetInt(cast.ToInt64(value))

	case reflect.Float32, reflect.Float64:
		field.SetFloat(cast.ToFloat64(value))

	case reflect.String:
		field.SetString(cast.ToString(value))

	case reflect.Bool:
		field.SetBool(cast.ToBool(value))

	case reflect.Ptr:
		if field.IsNil() {
			field.Set(reflect.New(field.Type().Elem()))
		}
		return p.setField(field.Elem(), fieldType.Elem(), value)

	default:
		return fmt.Errorf("unsupported protobuf field type: %v", fieldType.Kind())
	}

	return nil
}

func (p *ProtoFieldSetter) setTimestampField(field reflect.Value, pgType string, value interface{}) error {
	var t time.Time
	var err error

	switch pgType {
	case "timestamp", "timestamp without time zone":
		t, err = p.converter.toTimestamp(value, false)
	case "timestamp with time zone", "timestamptz":
		t, err = p.converter.toTimestamp(value, true)
	case "date":
		t, err = p.converter.toDate(value)
	default:
		return fmt.Errorf("unsupported timestamp type: %s", pgType)
	}

	if err != nil {
		return fmt.Errorf("converting to timestamp: %v", err)
	}

	if field.IsNil() {
		field.Set(reflect.ValueOf(timestamppb.New(t)))
	} else {
		field.Interface().(*timestamppb.Timestamp).Seconds = t.Unix()
		field.Interface().(*timestamppb.Timestamp).Nanos = int32(t.Nanosecond())
	}

	return nil
}

func (p *ProtoFieldSetter) setStructField(field reflect.Value, pgType string, value interface{}) error {
	switch pgType {
	case "jsonb", "json":
		// For JSON types
		var mapValue map[string]interface{}

		switch v := value.(type) {
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
			return fmt.Errorf("unexpected type for Struct: %T", value)
		}

		// Convert to structpb.Struct
		pbStruct, err := structpb.NewStruct(mapValue)
		if err != nil {
			return fmt.Errorf("converting to protobuf Struct: %v", err)
		}

		// Set the field value
		field.Set(reflect.ValueOf(pbStruct))
		return nil

	default:
		return fmt.Errorf("unsupported type for Struct conversion: %s", pgType)
	}
}
