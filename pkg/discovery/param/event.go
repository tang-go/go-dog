package param

const (
	//Listen 监听事件
	Listen int8 = iota
	//Heart 心跳
	Heart
	//Reg 注册
	Reg
)

//Event 事件
type Event struct {
	Codec
	Cmd   int8
	Label string
	Data  *Data
}

//All 获取所有
type All struct {
	Codec
	Label string
	Datas []*Data
}
