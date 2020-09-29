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

	//GetClient 获取客户端
	GetClient() Client

	//SetServiceFlowLimit 设置服务端最大流量限制
	SetServiceFlowLimit(max int64)

	//SetClientFlowLimit 设置客户端最大流量限制
	SetClientFlowLimit(max int64)

	//RegisterRPC 注册RPC方法
	RegisterRPC(name string, level int8, isAuth bool, explain string, fn interface{})

	//RegisterAPI 注册API方法--注册给网管
	RegisterAPI(methodname, version, path string, kind HTTPKind, level int8, isAuth bool, explain string, fn interface{})

	//GetCfg 获取配置
	GetCfg() Cfg

	//Run 启动服务
	Run() error
}
