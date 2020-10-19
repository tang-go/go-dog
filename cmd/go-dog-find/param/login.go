package param

const (
	//DisType 服务发现方
	DisType int8 = iota
	//RegType 服务注册方
	RegType
)

//LoginReq 登陆请求
type LoginReq struct {
	Codec
	Type int8
	Data *Data
}

//Data 数据
type Data struct {
	Label string
	Key   string
	Value string
}
