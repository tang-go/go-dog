package plugins

//Router RPC路由
type Router interface {

	//RegisterByMethod 注册方法
	RegisterByMethod(name string, fn interface{}) (arg map[string]interface{}, reply map[string]interface{})

	//GetMethodArg 获取方法请求的参数
	GetMethodArg(method string) (interface{}, bool)

	//Call 调用方法
	Call(ctx Context, method string, arg interface{}) (interface{}, error)
}
