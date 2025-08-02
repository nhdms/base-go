package common

import (
	"fmt"
	"net/http"
	neturl "net/url"
	"reflect"
	"strconv"
	"strings"

	"github.com/goccy/go-json"
)

type HTTPResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Code    int         `json:"code,omitempty"`
	Data    interface{} `json:"data"`
}

func NewSuccessHTTPResponse(data any) *HTTPResponse {
	return &HTTPResponse{Success: true, Data: data}
}

func NewErrorHTTPResponse(message string) *HTTPResponse {
	return &HTTPResponse{Success: false, Message: message}
}

func NewErrorCodeHTTPResponse(code int) *HTTPResponse {
	return &HTTPResponse{Success: false, Code: code, Message: Code2Message[code]}
}

func NewErrorWithCodeHTTPResponse(code int, message string) *HTTPResponse {
	return &HTTPResponse{Success: false, Code: code, Message: message}
}

func NewErrorWithDataHTTPResponse(code int, message string, data any) *HTTPResponse {
	return &HTTPResponse{Success: false, Code: code, Message: message, Data: data}
}

func (h HTTPResponse) ToJSON() string {
	v, _ := json.Marshal(h)
	return string(v)
}

func ParseSortParams(p string) ([]string, error) {
	if p == "" {
		return nil, nil
	}

	fields := strings.Split(p, ",")
	var orderBys []string
	for _, field := range fields {
		parts := strings.Split(field, ":")
		if len(parts) != 2 || (parts[1] != "asc" && parts[1] != "desc") {
			return nil, fmt.Errorf("invalid sort format: %s", field)
		}
		orderBys = append(orderBys, fmt.Sprintf("%s %s", parts[0], strings.ToUpper(parts[1])))
	}
	return orderBys, nil
}

// ExtractHeader extracts a header from the request and set extracted value into the provided destination
// Supported types: string, int family types, uint familiy types, float family types, bool,
// []string, []int, []int64, []float64
func ExtractHeader(r *http.Request, key string, dest interface{}) error {
	values, ok := r.Header[key]
	if !ok || len(values) == 0 {
		return fmt.Errorf("header %s not found", key)
	}

	val := values[0]
	v := reflect.ValueOf(dest)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("destination must be a non-nil pointer")
	}

	switch v := v.Elem(); v.Kind() {
	case reflect.String:
		v.SetString(val)
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		i, err := strconv.ParseInt(val, 10, 64)
		if err != nil {
			return err
		}
		v.SetInt(i)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		i, err := strconv.ParseUint(val, 10, 64)
		if err != nil {
			return err
		}
		v.SetUint(i)
	case reflect.Float32, reflect.Float64:
		f, err := strconv.ParseFloat(val, 64)
		if err != nil {
			return err
		}
		v.SetFloat(f)
	case reflect.Bool:
		b, err := strconv.ParseBool(val)
		if err != nil {
			return err
		}
		v.SetBool(b)
	case reflect.Slice:
		elem := v.Elem()
		switch elem.Type().Elem().Kind() {
		case reflect.String:
			elem.Set(reflect.ValueOf(values))

		case reflect.Int:
			ints := make([]int, 0, len(values))
			for _, s := range values {
				i, err := strconv.Atoi(s)
				if err != nil {
					return err
				}
				ints = append(ints, i)
			}
			elem.Set(reflect.ValueOf(ints))

		case reflect.Int64:
			ints := make([]int64, 0, len(values))
			for _, s := range values {
				i, err := strconv.ParseInt(s, 10, 64)
				if err != nil {
					return err
				}
				ints = append(ints, i)
			}
			elem.Set(reflect.ValueOf(ints))

		case reflect.Float64:
			floats := make([]float64, 0, len(values))
			for _, s := range values {
				f, err := strconv.ParseFloat(s, 64)
				if err != nil {
					return err
				}
				floats = append(floats, f)
			}
			elem.Set(reflect.ValueOf(floats))

		default:
			return fmt.Errorf("unsupported slice element type: %s", elem.Type().Elem().Kind())
		}
	default:
		return fmt.Errorf("unsupported type: %s", v.Kind())
	}
	return nil
}

// ExtractQueryParams extracts a query parameter from the query string and returns its value
func ExtractQueryParams(queryParams, target string) (string, error) {
	values, err := neturl.ParseQuery(queryParams)
	if err != nil {
		return "", fmt.Errorf("failed to parse query params: %s", err.Error())
	}

	targetValues, ok := values[target]
	if !ok {
		return "", fmt.Errorf("target %s not found in query params", target)
	}

	return targetValues[0], nil
}
