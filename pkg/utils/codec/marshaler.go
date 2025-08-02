package codec

import (
	"github.com/goccy/go-json"
	"github.com/mitchellh/mapstructure"
)

func JsonMarshal(v interface{}, forceEmit bool) ([]byte, error) {
	if !forceEmit {
		return json.Marshal(v)
	}

	var m map[string]interface{}
	_ = mapstructure.Decode(v, &m)
	return json.Marshal(m)
}
