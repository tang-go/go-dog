package param

//LoginType 登陆类型
type LoginType int8

const (
	//DisType 服务发现方
	DisType LoginType = iota
	//RegType 服务注册方
	RegType
)

//LoginReq 登陆请求
type LoginReq struct {
	Codec
	//Type 类型
	Type LoginType
}

//标签
type Label int8

const (
	//APILabel api label
	APILabel Label = iota
	//RPCLabel rpc label
	RPCLabel
)

//RegReq 注册请求
type RegReq struct {
	Codec
	Label Label
	Data  Data
}

//ListenReq 监听请求
type ListenReq struct {
	Codec
	Label Label
}

//ListenReq 监听请求
type ListenRes struct {
	Codec
	Label Label
	Data  []Data
}

//Data 数据
type Data struct {
	Key   string
	Value string
}
