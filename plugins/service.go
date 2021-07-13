package plugins

import "github.com/tang-go/go-dog/metrics"

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

//HTTP HTTP接口
type HTTP interface {

	//Group APi组
	Group(group string) HTTP

	//Auth APi需要验证
	Auth() HTTP

	//NoAuth APi需要不验证
	NoAuth() HTTP

	//Version APi版本
	Version(version string) HTTP

	//Auth APi等级
	Level(level int8) HTTP

	//Class 对象
	Class(class string) HTTP

	//GET APi GET路由
	GET(method string, path string, explain string, fn interface{})

	//POST POST路由
	POST(method string, path string, explain string, fn interface{})

	//PUT PUT路由
	PUT(method string, path string, explain string, fn interface{})

	//DELETE DELETE路由
	DELETE(method string, path string, explain string, fn interface{})
}

//RPC RPC接口
type RPC interface {

	//Auth 需要验证
	Auth() RPC

	//NoAuth 需要不验证
	NoAuth() RPC

	//Class 对象
	Class(class string) RPC

	//Method 方法
	Method(method string, explain string, fn interface{})
}

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

	//RPC 创建rpc
	RPC() RPC

	//HTTP 创建http
	HTTP(gate string) HTTP

	//APIRegIntercept API注册拦截器
	APIRegIntercept(f func(gate, group, url string, level int8, isAuth bool, explain string))

	//Run 启动服务
	Run() error

	//AddMetricValue 添加metric采集的值
	AddMetricValue(metricValue []*metrics.MetricValue)
}
