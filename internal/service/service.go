package service

import (
	"fmt"
	customerror "go-dog/error"
	"go-dog/header"
	"go-dog/internal/client"
	"go-dog/internal/codec"
	"go-dog/internal/config"
	"go-dog/internal/context"
	"go-dog/internal/interceptor"
	"go-dog/internal/limit"
	"go-dog/internal/register"
	"go-dog/internal/router"
	"go-dog/internal/rpc"
	"go-dog/pkg/log"
	"go-dog/pkg/recover"
	"go-dog/plugins"
	"go-dog/serviceinfo"
	"net"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

const (
	_MaxServiceRequestCount = 100000
)

//Service 服务
type Service struct {
	//服务名称
	name string
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
func CreateService(ttl int64, param ...interface{}) plugins.Service {
	service := &Service{
		close: 0,
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
		service.register = register.NewEtcdRegister(service.cfg.GetEtcd(), ttl)
	}
	if service.router == nil {
		//默认路由
		service.router = router.NewRouter(service.codec)
	}
	if service.limit == nil {
		//默认限流插件
		service.limit = limit.NewLimit(_MaxServiceRequestCount)
	}
	if service.interceptor == nil {
		//默认拦截器
		service.interceptor = interceptor.NewInterceptor()
	}
	if service.client == nil {
		//默认客户端
		service.client = client.NewClient(ttl, service.cfg)
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

//SetName 设置服务名称
func (s *Service) SetName(name string) {
	s.name = name
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
	log.Traceln("注册RPC方法:", method.Name, "说明:", method.Explain)
}

//POST POST方法
func (s *Service) POST(methodname, version, path string, level int8, isAuth bool, explain string, fn interface{}) {
	s.RegisterAPI(methodname, version, path, plugins.POST, level, isAuth, explain, fn)
}

//GET GET方法
func (s *Service) GET(methodname, version, path string, level int8, isAuth bool, explain string, fn interface{}) {
	s.RegisterAPI(methodname, version, path, plugins.GET, level, isAuth, explain, fn)
}

//RegisterAPI 注册API方法--注册给网管
func (s *Service) RegisterAPI(methodname, version, path string, kind plugins.HTTPKind, level int8, isAuth bool, explain string, fn interface{}) {
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
	log.Tracef("注册API接口:%s,路由:api/%s/%s/%s", api.Name, s.name, api.Version, api.Path)
}

//Run 启动服务
func (s *Service) Run() error {
	c := make(chan os.Signal)
	//监听指定信号
	signal.Notify(c, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		err := s.run()
		if err != nil {
			panic(err.Error())
		}
	}()
	msg := <-c
	atomic.AddInt32(&s.close, 1)
	s.Close()
	return fmt.Errorf("收到kill信号:%s", msg)
}

func (s *Service) run() error {
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
		go s.serveConn(conn)
	}
}

// ServeConn 拦截一个链接
func (s *Service) serveConn(conn net.Conn) {
	serviceRPC := rpc.NewServiceRPC(conn)
	serviceRPC.RegisterCallNotice(
		func(req *header.Request) *header.Response {
			defer recover.Recover()
			//服务器关闭了 直接关闭
			if atomic.LoadInt32(&s.close) > 0 {
				rep := new(header.Response)
				rep.ID = req.ID
				rep.Method = req.Method
				rep.Name = req.Name
				rep.Reply = nil
				rep.Error = customerror.EnCodeError(customerror.InternalServerError, "服务器关闭")
				return rep
			}
			//此处等待处理进程处理
			s.wait.Add(1)
			defer s.wait.Done()

			if s.limit.IsLimit() {
				rep := new(header.Response)
				rep.ID = req.ID
				rep.Method = req.Method
				rep.Name = req.Name
				rep.Reply = nil
				rep.Error = customerror.EnCodeError(customerror.SeviceLimitError, "超过服务每秒限制流量")
				return rep
			}

			rep := new(header.Response)
			rep.ID = req.ID
			rep.Method = req.Method
			rep.Name = req.Name
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
			ctx.SetClient(s.client)
			for key, value := range req.Data {
				ctx.SetData(key, value)
			}
			ctx = context.WithTimeout(ctx, ttl)
			if s.interceptor != nil {
				s.interceptor.Request(ctx, req.Name, req.Method, req.Arg)
			}
			back, err := s.router.Call(ctx, req)
			if s.interceptor != nil {
				s.interceptor.Respone(ctx, req.Name, req.Method, back)
			}
			if err != nil {
				rep.Error = customerror.DeCodeError(err)
				return rep
			}
			rep.Reply = back
			return rep
		})
}

//Close 关闭服务
func (s *Service) Close() {
	s.wait.Wait()
	s.register.Cancellation()
	s.limit.Close()
	s.client.Close()
}
