package pg_converter

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/cast"
	"time"
)

var (
	SupportTimeFormats = []string{
		"2006-01-02 15:04:05.999999999Z07:00",
		"2006-01-02 15:04:05.999999999",
		"2006-01-02 15:04:05.999999Z07:00",
		"2006-01-02 15:04:05.999999",
		"2006-01-02 15:04:05.999999-07",
		"2006-01-02 15:04:05Z07:00",
		"2006-01-02 15:04:05",
		"2006-01-02T15:04:05.999999999Z07:00",
		"2006-01-02T15:04:05.999999999",
		"2006-01-02T15:04:05.999999Z07:00",
		"2006-01-02T15:04:05.999999",
		"2006-01-02T15:04:05Z07:00",
		"2006-01-02T15:04:05",
	}
)

// PostgreSQLTypeConverter handles conversion from WAL message types to Go types
type PostgreSQLTypeConverter struct {
	location *time.Location
}

// NewPostgreSQLTypeConverter creates a new type converter
func NewPostgreSQLTypeConverter(location *time.Location) *PostgreSQLTypeConverter {
	if location == nil {
		location = time.UTC
	}
	return &PostgreSQLTypeConverter{
		location: location,
	}
}

// ConvertValue converts a WAL column value to the appropriate Go type
func (c *PostgreSQLTypeConverter) ConvertValue(pgType string, value interface{}) (interface{}, error) {
	if value == nil {
		return nil, nil
	}

	// Handle array types
	if len(pgType) > 2 && pgType[0] == '_' {
		return c.convertArray(pgType[1:], value)
	}

	switch pgType {
	case "smallint", "int2":
		return cast.ToInt16E(value)
	case "integer", "int", "int4":
		return cast.ToInt32E(value)
	case "bigint", "int8":
		return cast.ToInt64E(value)
	case "real", "float4":
		return cast.ToFloat32E(value)
	case "double precision", "float8":
		return cast.ToFloat64E(value)
	case "character varying", "varchar", "character", "char", "text":
		return cast.ToStringE(value)
	case "boolean", "bool":
		return cast.ToBoolE(value)
	case "timestamp", "timestamp without time zone":
		return c.toTimestamp(value, false)
	case "timestamp with time zone", "timestamptz":
		return c.toTimestamp(value, true)
	case "date":
		return c.toDate(value)
	case "json", "jsonb":
		return c.toJSON(value)
	case "inet", "cidr", "uuid":
		return cast.ToStringE(value)
	default:
		if value == nil {
			return nil, nil
		}
		return fmt.Sprintf("%v", value), nil
	}
}

func (c *PostgreSQLTypeConverter) toTimestamp(v interface{}, withTZ bool) (time.Time, error) {
	switch val := v.(type) {
	case string:
		var t time.Time
		var err error
		for _, format := range SupportTimeFormats {
			t, err = time.Parse(format, val)
			if err == nil {
				if withTZ {
					return t, nil
				}
				return t.In(c.location), nil
			}
		}
		return time.Time{}, fmt.Errorf("unable to parse timestamp: %v", val)
	case time.Time:
		if withTZ {
			return val, nil
		}
		return val.In(c.location), nil
	case nil:
		return time.Time{}, nil
	default:
		return time.Time{}, fmt.Errorf("unexpected type for timestamp: %T", v)
	}
}

func (c *PostgreSQLTypeConverter) toDate(v interface{}) (time.Time, error) {
	switch val := v.(type) {
	case string:
		t, err := time.Parse("2006-01-02", val)
		if err != nil {
			// Try parsing as timestamp and extract date
			if t, err := c.toTimestamp(val, false); err == nil {
				return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, c.location), nil
			}
			return time.Time{}, fmt.Errorf("parsing date: %v", err)
		}
		return t, nil
	case time.Time:
		return time.Date(val.Year(), val.Month(), val.Day(), 0, 0, 0, 0, c.location), nil
	case nil:
		return time.Time{}, nil
	default:
		return time.Time{}, fmt.Errorf("unexpected type for date: %T", v)
	}
}

func (c *PostgreSQLTypeConverter) toJSON(v interface{}) (interface{}, error) {
	switch val := v.(type) {
	case string:
		var result interface{}
		if err := json.Unmarshal([]byte(val), &result); err != nil {
			return nil, fmt.Errorf("parsing JSON: %v", err)
		}
		return result, nil
	case []byte:
		var result interface{}
		if err := json.Unmarshal(val, &result); err != nil {
			return nil, fmt.Errorf("parsing JSON: %v", err)
		}
		return result, nil
	case nil:
		return nil, nil
	default:
		return v, nil
	}
}

func (c *PostgreSQLTypeConverter) convertArray(elemType string, v interface{}) (interface{}, error) {
	switch val := v.(type) {
	case string:
		var result []interface{}
		if err := json.Unmarshal([]byte(val), &result); err != nil {
			return nil, fmt.Errorf("parsing array: %v", err)
		}

		for i, elem := range result {
			converted, err := c.ConvertValue(elemType, elem)
			if err != nil {
				return nil, fmt.Errorf("converting array element %d: %v", i, err)
			}
			result[i] = converted
		}
		return result, nil
	case nil:
		return nil, nil
	default:
		return nil, fmt.Errorf("array value must be string, got %T", v)
	}
}
