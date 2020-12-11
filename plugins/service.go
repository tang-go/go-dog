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

	//GetCodec 获取编码插件
	GetCodec() Codec

	//GetCfg 获取配置
	GetCfg() Cfg

	//GetLimit 获取限流插件
	GetLimit() Limit

	//GetClient 获取客户端
	GetClient() Client

	//Auth 验证函数
	Auth(fun func(ctx Context, method, token string) error)

	//RegisterRPC 	注册RPC方法
	RPC(name string, level int8, isAuth bool, explain string, fn interface{})

	//POST 			注册POST方法
	POST(methodname, version, path string, level int8, isAuth bool, explain string, fn interface{})

	//GET GET方法
	GET(methodname, version, path string, level int8, isAuth bool, explain string, fn interface{})

	//HTTP 创建http
	HTTP() API

	//APIRegIntercept API注册拦截器
	APIRegIntercept(f func(url, explain string))

	//Run 启动服务
	Run() error
}
