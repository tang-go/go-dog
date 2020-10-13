package plugins

//HTTPKind http请求类型
type HTTPKind string

const (
	//GET  get请求
	GET HTTPKind = "GET"
	//POST post请求
	POST HTTPKind = "POST"
	//PUT put请求
	PUT HTTPKind = "PUT"
	//DELETE delete请求
	DELETE HTTPKind = "DELETE"
)

//Service 服务接口
type Service interface {

	//SetName 设置服务名称
	SetName(name string)

	//GetCodec 获取编码插件
	GetCodec() Codec

	//GetCfg 获取配置
	GetCfg() Cfg

	//GetLimit 获取限流插件
	GetLimit() Limit

	//GetClient 获取客户端
	GetClient() Client

	//RegisterRPC 	注册RPC方法
	//name			方法名称
	//level			方法等级
	//isAuth		是否需要鉴权
	//explain		方法说明
	//fn			注册的方法
	RPC(name string, level int8, isAuth bool, explain string, fn interface{})

	//POST 			注册POST方法
	//methodname 	API方法名称
	//version 		API方法版本
	//path 			API路由
	//level 		API等级
	//isAuth 		是否需要鉴权
	//explain		方法描述
	//fn 			注册的方法
	POST(methodname, version, path string, level int8, isAuth bool, explain string, fn interface{})

	//GET GET方法
	//methodname 	API方法名称
	//version 		API方法版本
	//path 			API路由
	//level 		API等级
	//isAuth 		是否需要鉴权
	//explain		方法描述
	//fn 			注册的方法
	GET(methodname, version, path string, level int8, isAuth bool, explain string, fn interface{})

	//Run 启动服务
	Run() error
}
