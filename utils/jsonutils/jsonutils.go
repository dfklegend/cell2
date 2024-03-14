package jsonutils

import jsoniter "github.com/json-iterator/go"

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func Marshal(obj any) string {
	bytes, err := json.Marshal(&obj)
	if err != nil {
		return ""
	}
	return string(bytes)
}

func Unmarshal(input []byte, obj any) {
	json.Unmarshal(input, &obj)
}
