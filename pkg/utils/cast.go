package utils

import (
	"github.com/goccy/go-json"
	"google.golang.org/protobuf/types/known/structpb"
)

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
	v, _ := structpb.NewStruct(m)
	return v
}
