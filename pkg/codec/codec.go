package codec

import (
	"bytes"
	"encoding/json"

	"github.com/vmihailenco/msgpack"
)

//Codec 编码器
type Codec struct {
}

//NewCodec 新建一个编码器
func NewCodec() *Codec {
	codec := new(Codec)
	return codec
}

//EnCode 编码
func (c *Codec) EnCode(code string, v interface{}) ([]byte, error) {
	switch code {
	case "json":
		return json.Marshal(v)
	default:
		return msgpack.Marshal(v)
	}

}

//DeCode 编码
func (c *Codec) DeCode(code string, buff []byte, v interface{}) error {
	switch code {
	case "json":
		d := json.NewDecoder(bytes.NewReader(buff))
		d.UseNumber()
		return d.Decode(v)
	default:
		return msgpack.Unmarshal(buff, v)
	}
}
