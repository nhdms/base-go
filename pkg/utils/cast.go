package utils

import (
	"fmt"
	"github.com/goccy/go-json"
	_struct "github.com/golang/protobuf/ptypes/struct"
	"github.com/nhdms/base-go/proto/exmsg/models"
	"github.com/spf13/cast"
	"google.golang.org/protobuf/types/known/structpb"
	"reflect"
	"strings"
)

// tagKey defines a structure tag name for ConvertStructPB.
const tagKey = "structpb"

func ToJSONString(v interface{}) string {
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(jsonBytes)
}

func ToJSONByte(v interface{}) []byte {
	jsonBytes, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	return jsonBytes
}

func Map2ProtoStruct(m map[string]interface{}) *structpb.Struct {
	v, b := structpb.NewStruct(m)
	fmt.Println(b)
	return v
}

// Convert converts a structpb.Struct object to a concrete object.
func Convert(src *structpb.Struct, dst interface{}) error {
	return convertStruct(src, reflect.ValueOf(dst))
}

func toPrimitive(src *structpb.Value) (reflect.Value, bool) {
	switch t := src.GetKind().(type) {
	case *structpb.Value_BoolValue:
		return reflect.ValueOf(t.BoolValue), true
	case *structpb.Value_NullValue:
		return reflect.ValueOf(nil), true
	case *structpb.Value_NumberValue:
		return reflect.ValueOf(t.NumberValue), true
	case *structpb.Value_StringValue:
		return reflect.ValueOf(t.StringValue), true
	default:
		return reflect.Value{}, false
	}
}

func convertValue(src *structpb.Value, dest reflect.Value) error {
	dst := reflect.Indirect(dest)
	if v, ok := toPrimitive(src); ok {
		if !v.Type().AssignableTo(dst.Type()) {
			if !v.Type().ConvertibleTo(dst.Type()) {
				return fmt.Errorf("cannot assign %T to %s", src.GetKind(), dst.Type())
			}
			v = v.Convert(dst.Type())
		}
		dst.Set(v)
		return nil
	}
	switch t := src.GetKind().(type) {
	case *structpb.Value_ListValue:
		return convertList(t.ListValue, dst)
	case *structpb.Value_StructValue:
		return convertStruct(t.StructValue, dst)
	default:
		return fmt.Errorf("unsuported value: %T", src.GetKind())
	}
}

func convertList(src *structpb.ListValue, dest reflect.Value) error {
	dst := reflect.Indirect(dest)
	if dst.Kind() != reflect.Slice {
		return fmt.Errorf("cannot convert %T to %s", src, dst.Type())
	}
	values := src.GetValues()
	elemType := dst.Type().Elem()
	converted := make([]reflect.Value, len(values))
	for i, value := range values {
		element := reflect.New(elemType).Elem()
		if err := convertValue(value, element); err != nil {
			return err
		}
		converted[i] = element
	}
	dst.Set(reflect.Append(dst, converted...))
	return nil
}

func convertStruct(src *structpb.Struct, dest reflect.Value) error {
	dst := reflect.Indirect(dest)
	if dst.Kind() == reflect.Struct {
		fields := src.GetFields()
		for i := 0; i < dst.NumField(); i++ {
			target := dst.Field(i)
			field := dst.Type().Field(i)
			name := field.Tag.Get(tagKey)
			if name == "" {
				name = strings.ToLower(field.Name)
			}
			if v, ok := fields[name]; ok {
				if err := convertValue(v, target); err != nil {
					return err
				}
			}
		}
		return nil
	} else if dst.Kind() == reflect.Map {
		elemType := dst.Type().Elem()
		mapType := reflect.MapOf(reflect.TypeOf(string("")), elemType)
		aMap := reflect.MakeMap(mapType)
		fields := src.GetFields()
		for key, value := range fields {
			element := reflect.New(elemType).Elem()
			if err := convertValue(value, element); err != nil {
				return err
			}
			aMap.SetMapIndex(reflect.ValueOf(key), element)
		}
		dst.Set(aMap)
		return nil
	}

	return fmt.Errorf("cannot convert %T to %s", src, dst.Type())
}

func ToArgs[T any](v T) *structpb.Value {
	rv := reflect.ValueOf(v)

	// Handle nil for interface and pointer types
	if (rv.Kind() == reflect.Interface || rv.Kind() == reflect.Ptr) && rv.IsNil() {
		return nil
	}

	switch reflect.ValueOf(v).Kind() {
	case reflect.String:
		return &structpb.Value{
			Kind: &structpb.Value_StringValue{
				StringValue: reflect.ValueOf(v).String(),
			},
		}
	case reflect.Bool:
		return &structpb.Value{
			Kind: &structpb.Value_BoolValue{
				BoolValue: reflect.ValueOf(v).Bool(),
			},
		}
	case reflect.Float32, reflect.Float64, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return &structpb.Value{
			Kind: &structpb.Value_NumberValue{
				NumberValue: cast.ToFloat64(reflect.ValueOf(v).Interface()),
			},
		}
	case reflect.Slice:
		// reflect for loop this slice and recursive each item
		sliceVal := reflect.ValueOf(v)
		values := make([]*structpb.Value, sliceVal.Len())

		for i := 0; i < sliceVal.Len(); i++ {
			item := sliceVal.Index(i).Interface()
			values[i] = ToArgs(item)
		}

		return &structpb.Value{
			Kind: &structpb.Value_ListValue{
				ListValue: &structpb.ListValue{
					Values: values,
				},
			},
		}
	default:
		return nil
	}
}

func GetValueFromArgs(args *_struct.Value) interface{} {
	switch args.Kind.(type) {
	case *_struct.Value_StringValue:
		return args.GetStringValue()
	case *_struct.Value_NumberValue:
		return args.GetNumberValue()
	case *_struct.Value_BoolValue:
		return args.GetBoolValue()
	case *_struct.Value_ListValue:
		list := make([]interface{}, 0)
		for _, v := range args.GetListValue().GetValues() {
			list = append(list, GetValueFromArgs(v))
		}
		return list
	default:
		return nil
	}
}

func ConvertSQLRowToDest(row *models.Row, i interface{}) error {
	if row == nil {
		return nil
	}

	return json.Unmarshal(row.Values, i)
}
