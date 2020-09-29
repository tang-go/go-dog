package header

import (
	"github.com/vmihailenco/msgpack"
)

//MsgPackCode MsgPack编码
type MsgPackCode struct {
}

//EnCode 编码
func (r *MsgPackCode) EnCode(inter interface{}) ([]byte, error) {
	return msgpack.Marshal(inter)
}

//DeCode 解码
func (r *MsgPackCode) DeCode(buff []byte, inter interface{}) error {
	return msgpack.Unmarshal(buff, inter)
}
