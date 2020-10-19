package redis

import "encoding/json"

//Result 结果
type Result struct {
	value string
}

//DeCode 编码
func (r *Result) DeCode(v interface{}) error {
	return json.Unmarshal([]byte(r.value), v)
}
