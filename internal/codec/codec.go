package codec

import "encoding/json"

//Codec 编码器
type Codec struct {
}

//NewCodec 新建一个编码器
func NewCodec() *Codec {
	codec := new(Codec)
	return codec
}

//EnCode 编码
func (c *Codec) EnCode(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

//DeCode 编码
func (c *Codec) DeCode(buff []byte, v interface{}) error {
	return json.Unmarshal(buff, v)
}
