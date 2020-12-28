<h1 align="center">go-dog微服务框架</h1>
<div align="center">
go-dog是一个可插拔式的微服务框架，包含服务注册与发现，负载均衡，链路追踪，服务降级，函数级别的熔断，服务限流，服务治理平台
</div>

### 使用

- 请求查看项目go-dog-example https://github.com/tang-go/go-dog-example

### 接口说明
- plugins/service.go 服务端接口
```go
//GetCodec 获取编码插件
GetCodec() Codec

//GetCfg 获取配置
GetCfg() Cfg

//GetLimit 获取限流插件
GetLimit() Limit

//GetClient 获取客户端
GetClient() Client

//Auth 验证函数
Auth(fun func(ctx Context, token string) error)

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
```

- plugins/cfg.go 配置文件接口
```go
//GetServerName 获取服务名称
GetServerName() string

//GetExplain 获取服务说明
GetExplain() string

//GetPort 获取端口
GetPort() int

//GetDiscovery 获取服务发现配置
GetDiscovery() []string

//GetRedis 获取redis配置
GetRedis() []string

//GetKafka 获取kfaka地址
GetKafka() []string

//GetNats 获取nats地址
GetNats() []string

//GetRocketMq 获取rocketmq地址
GetRocketMq() []string

//GetNsq 获取nsq地址
GetNsq() []string

//GetReadMysql 获取ReadMysql地址
GetReadMysql() *config.MysqlCfg

//GetWriteMysql 获取GetWriteMysql地址
GetWriteMysql() *config.MysqlCfg

//GetHost 获取本机地址配置
GetHost() string

//GetRunmode 获取runmode地址配置
GetRunmode() string
```

- plugins/client.go 客户端接口
```go
//GetCodec 获取编码插件
GetCodec() Codec

//GetCfg 获取配置
GetCfg() Cfg

//GetDiscovery 获取服务发现
GetDiscovery() Discovery

//GetFusing 获取熔断插件
GetFusing() Fusing

//GetLimit 获取限流插件
GetLimit() Limit

//GetAllService 获取所有服务
GetAllService() (services []*serviceinfo.ServiceInfo)

//Call 调用函数
Call(ctx Context, mode Mode, name string, method string, args interface{}, reply interface{}) error

//SendRequest 发生请求
SendRequest(ctx Context, mode Mode, name string, method string, code string, args []byte) (reply []byte, e error)

//Close 关闭
Close()
```
### 其他

- go-dog-tool 工具 https://github.com/tang-go/go-dog-tool
- go-dog-vue 前端 https://github.com/tang-go/go-dog-vue
- go-dog-example 例子 https://github.com/tang-go/go-dog-example

### trace部署
docker run -d --name jaeger -e COLLECTOR_ZIPKIN_HTTP_PORT=9411 -p 5775:5775/udp -p 6831:6831/udp -p 6832:6832/udp -p 5778:5778 -p 16686:16686 -p 14268:14268 -p 9411:9411 jaegertracing/all-in-one:1.12

