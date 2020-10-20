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
	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/pkg/client"
	"github.com/tang-go/go-dog/pkg/codec"
	"github.com/tang-go/go-dog/pkg/config"
	"github.com/tang-go/go-dog/pkg/context"
	"github.com/tang-go/go-dog/pkg/limit"
	"github.com/tang-go/go-dog/pkg/register"
	"github.com/tang-go/go-dog/pkg/router"
	"github.com/tang-go/go-dog/pkg/rpc"
	"github.com/tang-go/go-dog/plugins"
	"github.com/tang-go/go-dog/recover"
	"github.com/tang-go/go-dog/serviceinfo"
)

const (
	_MaxServiceRequestCount = 100000
)

//Service 服务
type Service struct {
	//服务名称
	name string
	//验证插件
	auth func(ctx plugins.Context, token string) error
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
	if service.client == nil {
		//默认客户端
		service.client = client.NewClient(service.cfg)
	}
	//初始化日志
	switch service.cfg.GetRunmode() {
	case "panic":
		log.SetLevel(log.PanicLevel)
		break
	case "fatal":
		log.SetLevel(log.FatalLevel)
		break
	case "error":
		log.SetLevel(log.ErrorLevel)
		break
	case "warn":
		log.SetLevel(log.WarnLevel)
		break
	case "info":
		log.SetLevel(log.InfoLevel)
		break
	case "debug":
		log.SetLevel(log.DebugLevel)
		break
	case "trace":
		log.SetLevel(log.TraceLevel)
		break
	default:
		log.SetLevel(log.TraceLevel)
		break
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

//POST POST方法
func (s *Service) POST(methodname, version, path string, level int8, isAuth bool, explain string, fn interface{}) {
	s._RegisterAPI(methodname, version, path, plugins.POST, level, isAuth, explain, fn)
}

//GET GET方法
func (s *Service) GET(methodname, version, path string, level int8, isAuth bool, explain string, fn interface{}) {
	s._RegisterAPI(methodname, version, path, plugins.GET, level, isAuth, explain, fn)
}

//RegisterAPI 注册API方法--注册给网管
func (s *Service) _RegisterAPI(methodname, version, path string, kind plugins.HTTPKind, level int8, isAuth bool, explain string, fn interface{}) {
	req, rep := s.router.RegisterByMethod(methodname, fn)
	api := &serviceinfo.API{
		Name:     methodname,
		Level:    level,
		Explain:  explain,
		IsAuth:   isAuth,
		Request:  req,
		Response: rep,
		Version:  version,
		Path:     path,
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
	log.Tracef("注册API接口:%s,路由:api/%s/%s/%s", api.Name, s.name, api.Version, api.Path)
}

//Auth 验证函数
func (s *Service) Auth(fun func(ctx plugins.Context, token string) error) {
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
	l, err := net.Listen("tcp", fmt.Sprintf("%s:%d", s.cfg.GetHost(), s.cfg.GetPort()))
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
		info := serviceinfo.ServiceInfo{
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
			log.Tracef("| %s | %s | %13v | %s | %s ",
				address,
				respone.Error.Error(),
				time.Now().Sub(start),
				name,
				method,
			)
		} else {
			log.Tracef("| %s | %s | %13v | %s | %s ",
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
			if s.GetCfg().GetRunmode() == "trace" {
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
			ctx := context.Background()
			ctx.SetAddress(req.Address)
			ctx.SetTraceID(req.TraceID)
			ctx.SetIsTest(req.IsTest)
			ctx.SetToken(req.Token)

			ctx.SetClient(s.client)
			for key, value := range req.Data {
				ctx.SetData(key, value)
			}
			ctx = context.WithTimeout(ctx, ttl)

			if argv, ok := s.router.GetMethodArg(req.Method); ok {
				err := s.codec.DeCode(req.Code, req.Arg, argv)
				if err != nil {
					rep.Error = customerror.EnCodeError(customerror.ParamError, "参数不合法")
					return rep
				}
				//先判断此方法是否需要鉴权
				if _, o := s.authMethod[strings.ToLower(req.Method)]; o {
					if s.auth != nil {
						if err := s.auth(ctx, req.Token); err != nil {
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
					rep.Error = customerror.EnCodeError(customerror.ParamError, "返回参数不合法")
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
}
