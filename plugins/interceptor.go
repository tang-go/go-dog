package plugins

//Interceptor 拦截器
type Interceptor interface {
	//Request 请求
	Request(ctx Context, servicename, method string, request interface{})

	//Respone 响应
	Respone(ctx Context, servicename, method string, respone interface{}, err error)
}
