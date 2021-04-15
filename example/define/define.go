package define

//服务定义
const (
	// Prefix 统一前缀
	Prefix = "example/"
)

//定义网关层
const (
	//ExampleGate 演示网关
	ExampleGate = Prefix + "examplegate"
)

//定义逻辑服务层
const (
	//ServiceOne 服务1
	ServiceOne = Prefix + "serviceone"
	//ServiceTwo 服务2
	ServiceTwo = Prefix + "servicetwo"
)
