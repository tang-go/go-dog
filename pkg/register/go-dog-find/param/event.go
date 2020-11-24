package param

const (
	Login int8 = iota
	//Listen 监听事件
	Listen
	//Heart 心跳
	Heart
	//Reg 注册
	Reg
)

//Event 事件
type Event struct {
	Codec
	Cmd  int8
	Data []byte
}
