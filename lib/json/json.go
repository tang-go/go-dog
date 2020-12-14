package json

import (
	"bytes"

	jsoniter "github.com/json-iterator/go"
)

//func Marshal(o interface{}) (string, error) {
//	b, e := json.Marshal(o)
//	return string(b), e
//}
//
//func Unmarshal(data []byte, v interface{}) error {
//	return json.Unmarshal(data, v)
//}
func Marshal(o interface{}) (string, error) {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	b, e := json.Marshal(&o)
	b = bytes.Replace(b, []byte("\\u003c"), []byte("<"), -1)
	b = bytes.Replace(b, []byte("\\u003e"), []byte(">"), -1)
	b = bytes.Replace(b, []byte("\\u0026"), []byte("&"), -1)

	return string(b), e
}

func Unmarshal(data []byte, v interface{}) error {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary
	return json.Unmarshal(data, v)
}
