package param

import "github.com/vmihailenco/msgpack"

//Codec 编码器
type Codec struct {
}

//EnCode 编码
func (c *Codec) EnCode(v interface{}) ([]byte, error) {
	return msgpack.Marshal(v)

}

//DeCode 编码
func (c *Codec) DeCode(buff []byte, v interface{}) error {
	return msgpack.Unmarshal(buff, v)
}
