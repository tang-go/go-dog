package service

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	customerror "github.com/tang-go/go-dog/error"
	"github.com/tang-go/go-dog/header"
	"github.com/tang-go/go-dog/jaeger"
	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/pkg/client"
	"github.com/tang-go/go-dog/pkg/codec"
	"github.com/tang-go/go-dog/pkg/config"
	"github.com/tang-go/go-dog/pkg/context"
	"github.com/tang-go/go-dog/pkg/limit"
	register "github.com/tang-go/go-dog/pkg/register/go-dog-find"
	"github.com/tang-go/go-dog/pkg/router"
	"github.com/tang-go/go-dog/pkg/rpc"
	"github.com/tang-go/go-dog/plugins"
	"github.com/tang-go/go-dog/recover"
	"github.com/tang-go/go-dog/serviceinfo"
)

const (
	_MaxServiceRequestCount = 100000
)

//API api路由组件
type API struct {
	api *serviceinfo.API
	s   *Service
}

func newAPI(s *Service, api *serviceinfo.API) plugins.API {
	return &API{
		s:   s,
		api: api,
	}
}

//APIGroup APi组
func (a *API) APIGroup(group string) plugins.API {
	a.api.Group = group
	return a
}

//APIAuth APi需要验证
func (a *API) APIAuth() plugins.API {
	a.api.IsAuth = true
	return a
}

//APINoAuth APi需要不验证
func (a *API) APINoAuth() plugins.API {
	a.api.IsAuth = false
	return a
}

//APIVersion APi版本
func (a *API) APIVersion(version string) plugins.API {
	a.api.Version = version
	return a
}

//APILevel APi等级
func (a *API) APILevel(level int8) plugins.API {
	a.api.Level = level
	return a
}

//GET APi GET路由
func (a *API) GET(name string, path string, explain string, fn interface{}) {
	a.api.Path = path
	if a.api.Group == "" {
		a.api.Group = a.s.name
	}
	if a.api.Level <= 0 {
		a.api.Level = 1
	}
	a.s._RegisterAPI(a.api.Gate, a.api.Group, name, a.api.Version, path, plugins.GET, a.api.Level, a.api.IsAuth, explain, fn)
}

//POST POST路由
func (a *API) POST(name string, path string, explain string, fn interface{}) {
	a.api.Path = path
	if a.api.Group == "" {
		a.api.Group = a.s.name
	}
	if a.api.Version == "" {
		a.api.Version = "v1"
	}
	if a.api.Level <= 0 {
		a.api.Level = 1
	}
	a.s._RegisterAPI(a.api.Gate, a.api.Group, name, a.api.Version, path, plugins.POST, a.api.Level, a.api.IsAuth, explain, fn)
}

//Service 服务
type Service struct {
	//服务名称
	name string
	//验证插件
	auth func(ctx plugins.Context, method, token string) error
	//配置插件
	cfg plugins.Cfg
	//注册中心插件
	register plugins.Register
	//路由插件
	router plugins.Router
	//限流插件
	limit plugins.Limit
	//链路追踪插件
	interceptor plugins.Interceptor
	//服务发现
	discovery plugins.Discovery
	//方法
	methods []*serviceinfo.Method
	//鉴权方法
	authMethod map[string]string
	//对网管注册的api
	api []*serviceinfo.API
	//客户端
	client plugins.Client
	//参数编码器
	codec plugins.Codec
	//退出信号
	close int32
	//api注册拦截器
	apiRegIntercept func(gate, group, url string, level int8, isAuth bool, explain string)
	//等待
	wait sync.WaitGroup
}

//CreateService 创建一个服务
func CreateService(name string, param ...interface{}) plugins.Service {
	service := &Service{
		close:      0,
		name:       name,
		authMethod: make(map[string]string),
	}
	for _, plugin := range param {
		if cfg, ok := plugin.(plugins.Cfg); ok {
			service.cfg = cfg
		}
		if register, ok := plugin.(plugins.Register); ok {
			service.register = register
		}
		if router, ok := plugin.(plugins.Router); ok {
			service.router = router
		}
		if limit, ok := plugin.(plugins.Limit); ok {
			service.limit = limit
		}
		if interceptor, ok := plugin.(plugins.Interceptor); ok {
			service.interceptor = interceptor
		}
		if codec, ok := plugin.(plugins.Codec); ok {
			service.codec = codec
		}
		if discovery, ok := plugin.(plugins.Discovery); ok {
			service.discovery = discovery
		}
		if client, ok := plugin.(plugins.Client); ok {
			service.client = client
		}
	}
	if service.cfg == nil {
		//默认配置
		service.cfg = config.NewConfig()
	}
	if service.codec == nil {
		//默认参数编码插件
		service.codec = codec.NewCodec()
	}
	if service.register == nil {
		//使用默认注册中心
		service.register = register.NewGoDogRegister(service.cfg.GetDiscovery())
	}
	if service.router == nil {
		//默认路由
		service.router = router.NewRouter()
	}
	if service.limit == nil {
		//默认限流插件
		service.limit = limit.NewLimit(_MaxServiceRequestCount)
	}
	if service.interceptor == nil {
		//链路追踪插件
		service.interceptor = jaeger.NewJaeger(name, service.cfg)
	}
	if service.client == nil {
		//默认客户端
		if service.discovery != nil {
			service.client = client.NewClient(service.cfg, service.discovery)
		} else {
			service.client = client.NewClient(service.cfg)
		}
	}
	return service
}

//GetClient 获取客户端
func (s *Service) GetClient() plugins.Client {
	return s.client
}

//GetCfg 获取配置
func (s *Service) GetCfg() plugins.Cfg {
	return s.cfg
}

//GetLimit 获取限流插件
func (s *Service) GetLimit() plugins.Limit {
	return s.limit
}

//GetCodec 获取编码插件
func (s *Service) GetCodec() plugins.Codec {
	return s.codec
}

//RPC 注册RPC方法
func (s *Service) RPC(name string, level int8, isAuth bool, explain string, fn interface{}) {
	req, rep := s.router.RegisterByMethod(name, fn)
	method := &serviceinfo.Method{
		Name:     name,
		Level:    level,
		Explain:  explain,
		IsAuth:   isAuth,
		Request:  req,
		Response: rep,
	}
	s.methods = append(s.methods, method)
	if isAuth {
		s.authMethod[strings.ToLower(name)] = name
	}
	log.Traceln("注册RPC方法:", method.Name, "说明:", method.Explain)
}

//HTTP 创建http
func (s *Service) HTTP(gate string) plugins.API {
	api := new(serviceinfo.API)
	api.Gate = gate
	return newAPI(s, api)
}

//APIRegIntercept API注册拦截器
func (s *Service) APIRegIntercept(f func(gate, group, url string, level int8, isAuth bool, explain string)) {
	s.apiRegIntercept = f
}

//RegisterAPI 注册API方法--注册给网管
func (s *Service) _RegisterAPI(gate, group, methodname, version, path string, kind plugins.HTTPKind, level int8, isAuth bool, explain string, fn interface{}) {
	req, rep := s.router.RegisterByMethod(methodname, fn)
	url := fmt.Sprintf("api/%s/%s/%s", s.name, version, path)
	api := &serviceinfo.API{
		Gate:     gate,
		Name:     methodname,
		Group:    group,
		Level:    level,
		Explain:  explain,
		IsAuth:   isAuth,
		Request:  req,
		Response: rep,
		Version:  version,
		Path:     url,
		Kind:     string(kind),
	}
	method := &serviceinfo.Method{
		Name:     methodname,
		Level:    level,
		Explain:  explain,
		IsAuth:   isAuth,
		Request:  req,
		Response: rep,
	}
	s.methods = append(s.methods, method)
	s.api = append(s.api, api)
	if isAuth {
		s.authMethod[strings.ToLower(methodname)] = methodname
	}
	if s.apiRegIntercept != nil {
		s.apiRegIntercept(gate, group, url, level, isAuth, explain)
	}
	log.Tracef("注册API接口:%s,路由:%s", api.Name, api.Path)
}

//Auth 验证函数
func (s *Service) Auth(fun func(ctx plugins.Context, method, token string) error) {
	s.auth = fun
}

//Run 启动服务
func (s *Service) Run() error {
	c := make(chan os.Signal)
	//监听指定信号
	signal.Notify(c, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		err := s._Run()
		if err != nil {
			panic(err.Error())
		}
	}()
	msg := <-c
	s.Close()
	return fmt.Errorf("收到kill信号:%s", msg)
}

//_Run 启动
func (s *Service) _Run() error {
	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", s.cfg.GetPort()))
	if err != nil {
		return err
	}
	defer l.Close()
	name := s.name
	if name == "" {
		name = s.cfg.GetServerName()
	}
	//注册RPC方法到etcd
	if len(s.methods) > 0 {
		info := serviceinfo.RPCServiceInfo{
			Name:    name,
			Address: s.cfg.GetHost(),
			Port:    s.cfg.GetPort(),
			Explain: s.cfg.GetExplain(),
			Methods: s.methods,
			Time:    time.Now().Format("2006-01-02 15:04:05"),
		}
		s.register.RegisterRPCService(context.Background(), &info)
	}
	//注册API方法到etcd
	if len(s.api) > 0 {
		info := serviceinfo.APIServiceInfo{
			Name:    name,
			Address: s.cfg.GetHost(),
			Port:    s.cfg.GetPort(),
			API:     s.api,
			Explain: s.cfg.GetExplain(),
			Time:    time.Now().Format("2006-01-02 15:04:05"),
		}
		s.register.RegisterAPIService(context.Background(), &info)
	}
	for {
		if atomic.LoadInt32(&s.close) > 0 {
			return nil
		}
		conn, err := l.Accept()
		if err != nil {
			log.Errorln(err.Error())
			continue
		}
		go s._ServeConn(conn)
	}
}

//_Log 日志
func (s *Service) _Log(address, name, method string, respone *header.Response) func() {
	start := time.Now()
	return func() {
		if respone.Error != nil {
			log.Infof("| %s | %s | %13v | %s | %s ",
				address,
				respone.Error.Error(),
				time.Now().Sub(start),
				name,
				method,
			)
		} else {
			log.Infof("| %s | %s | %13v | %s | %s ",
				address,
				"成功",
				time.Now().Sub(start),
				name,
				method,
			)
		}
	}
}

// ServeConn 拦截一个链接
func (s *Service) _ServeConn(conn net.Conn) {
	serviceRPC := rpc.NewServiceRPC(conn, s.codec)
	serviceRPC.RegisterCallNotice(
		func(req *header.Request) *header.Response {
			defer recover.Recover()
			rep := new(header.Response)
			rep.ID = req.ID
			rep.Method = req.Method
			rep.Name = req.Name
			rep.Code = req.Code
			if s.GetCfg().GetRunmode() == "trace" || s.GetCfg().GetRunmode() == "debug" || s.GetCfg().GetRunmode() == "info" {
				defer s._Log(req.Address, req.Name, req.Method, rep)()
			}
			//服务器关闭了 直接关闭
			if atomic.LoadInt32(&s.close) > 0 {
				rep.Error = customerror.EnCodeError(customerror.InternalServerError, "服务器关闭")
				return rep
			}
			//此处等待处理进程处理
			s.wait.Add(1)
			defer s.wait.Done()

			if s.limit.IsLimit() {
				rep.Error = customerror.EnCodeError(customerror.SeviceLimitError, "超过服务每秒限制流量")
				return rep
			}
			now := time.Now().UnixNano()
			ttl := req.TimeOut - now
			if ttl < 0 {
				//超时
				rep.Error = customerror.EnCodeError(customerror.RequestTimeout, "请求超时")
				return rep
			}
			//创建ctx
			datas := make(map[string][]byte)
			for key, value := range req.Data {
				datas[key] = value
			}
			ctx := context.NewContextByData(datas)
			ctx.SetAddress(req.Address)
			ctx.SetTraceID(req.TraceID)
			ctx.SetIsTest(req.IsTest)
			ctx.SetToken(req.Token)
			ctx.SetSource(req.Source)
			ctx.SetURL(req.URL)
			ctx.SetClient(s.client)

			ctx = context.WithTimeout(ctx, ttl)

			if argv, ok := s.router.GetMethodArg(req.Method); ok {
				err := s.codec.DeCode(req.Code, req.Arg, argv)
				if err != nil {
					rep.Error = customerror.EnCodeError(customerror.ParamError, "请求参数错误:"+err.Error())
					return rep
				}
				//先判断此方法是否需要鉴权
				if _, o := s.authMethod[strings.ToLower(req.Method)]; o {
					if s.auth != nil {
						if err := s.auth(ctx, req.Name, req.Token); err != nil {
							rep.Error = customerror.DeCodeError(err)
							return rep
						}
					}
				}
				if s.interceptor != nil {
					s.interceptor.Request(ctx, req.Name, req.Method, argv)
				}
				back, err := s.router.Call(ctx, req.Method, argv)
				if s.interceptor != nil {
					s.interceptor.Respone(ctx, rep.Name, rep.Method, back, err)
				}
				if err != nil {
					rep.Error = customerror.DeCodeError(err)
					return rep
				}
				reply, err := s.codec.EnCode(req.Code, back)
				if err != nil {
					rep.Error = customerror.EnCodeError(customerror.ParamError, "返回参数"+err.Error())
					return rep
				}
				rep.Reply = reply
				return rep
			}
			rep.Error = customerror.EnCodeError(customerror.RPCNotFind, "方法不存在")
			return rep
		})
}

//Close 关闭服务
func (s *Service) Close() {
	atomic.AddInt32(&s.close, 1)
	s.wait.Wait()
	s.register.Cancellation()
	s.limit.Close()
	s.client.Close()
	s.interceptor.Close()
}
