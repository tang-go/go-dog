package plugins

//Fusing 客户端熔断插件
type Fusing interface {
	//AddErrorMethod 添加请求发生错误的方法
	AddErrorMethod(servicekey, methodname string, err error)

	//AddMethod 添加请求
	AddMethod(servicekey, methodname string)

	//OpenFusing 设置某个服务方法强行开启熔断
	OpenFusing(servicekey, methodname string)

	//CloseFusing 设置某个服务方法关闭熔断
	CloseFusing(servicekey, methodname string)

	//IsFusing 是否熔断
	IsFusing(servicekey, methodname string) bool

	//Close 关闭
	Close()
}
