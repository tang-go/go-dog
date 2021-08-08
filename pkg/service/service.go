package service

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	customerror "github.com/tang-go/go-dog/error"
	"github.com/tang-go/go-dog/header"
	"github.com/tang-go/go-dog/jaeger"
	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/metrics"
	"github.com/tang-go/go-dog/pkg/client"
	"github.com/tang-go/go-dog/pkg/codec"
	"github.com/tang-go/go-dog/pkg/config"
	"github.com/tang-go/go-dog/pkg/context"
	"github.com/tang-go/go-dog/pkg/limit"
	consulRegister "github.com/tang-go/go-dog/pkg/register/consul"
	nacosRegister "github.com/tang-go/go-dog/pkg/register/nacos"
	"github.com/tang-go/go-dog/pkg/router"
	"github.com/tang-go/go-dog/pkg/rpc"
	"github.com/tang-go/go-dog/plugins"
	"github.com/tang-go/go-dog/recover"
	"github.com/tang-go/go-dog/serviceinfo"
)

type MetricOpts struct {
	NameSpace     string                 // 必填
	SystemName    string                 // 必填
	MetricsValues []*metrics.MetricValue // 必填
}

//HTTP api路由组件
type HTTP struct {
	api   *serviceinfo.API
	s     *Service
	class string
}

func newHTTP(s *Service, api *serviceinfo.API) plugins.HTTP {
	return &HTTP{
		s:   s,
		api: api,
	}
}

//Class 对象
func (a *HTTP) Class(class string) plugins.HTTP {
	a.class = class
	return a
}

//APIGroup APi组
func (a *HTTP) Group(group string) plugins.HTTP {
	a.api.Group = group
	return a
}

//Auth APi需要验证
func (a *HTTP) Auth() plugins.HTTP {
	a.api.IsAuth = true
	return a
}

//NoAuth APi需要不验证
func (a *HTTP) NoAuth() plugins.HTTP {
	a.api.IsAuth = false
	return a
}

//Version APi版本
func (a *HTTP) Version(version string) plugins.HTTP {
	a.api.Version = version
	return a
}

//Level APi等级
func (a *HTTP) Level(level int8) plugins.HTTP {
	a.api.Level = level
	return a
}

//GET APi GET路由
func (a *HTTP) GET(method string, path string, explain string, fn interface{}) {
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
	if a.class != "" {
		method = a.class + "." + method
	}
	a.s.RegisterAPI(a.api.Gate, a.api.Group, method, a.api.Version, path, plugins.GET, a.api.Level, a.api.IsAuth, explain, fn)
}

//DELETE APi DELETE路由
func (a *HTTP) DELETE(method string, path string, explain string, fn interface{}) {
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
	if a.class != "" {
		method = a.class + "." + method
	}
	a.s.RegisterAPI(a.api.Gate, a.api.Group, method, a.api.Version, path, plugins.DELETE, a.api.Level, a.api.IsAuth, explain, fn)
}

//POST POST路由
func (a *HTTP) POST(method string, path string, explain string, fn interface{}) {
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
	if a.class != "" {
		method = a.class + "." + method
	}
	a.s.RegisterAPI(a.api.Gate, a.api.Group, method, a.api.Version, path, plugins.POST, a.api.Level, a.api.IsAuth, explain, fn)
}

//PUT PUT路由
func (a *HTTP) PUT(method string, path string, explain string, fn interface{}) {
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
	if a.class != "" {
		method = a.class + "." + method
	}
	a.s.RegisterAPI(a.api.Gate, a.api.Group, method, a.api.Version, path, plugins.PUT, a.api.Level, a.api.IsAuth, explain, fn)
}

//RPC RPC注册
type RPC struct {
	method *serviceinfo.Method
	s      *Service
	class  string
}

func newRPC(s *Service, method *serviceinfo.Method) plugins.RPC {
	return &RPC{
		s:      s,
		method: method,
	}
}

//Class 对象
func (a *RPC) Class(class string) plugins.RPC {
	a.class = class
	return a
}

//Auth 需要验证
func (a *RPC) Auth() plugins.RPC {
	a.method.IsAuth = true
	return a
}

//NoAuth 需要不验证
func (a *RPC) NoAuth() plugins.RPC {
	a.method.IsAuth = false
	return a
}

//Level 等级
func (a *RPC) Level(level int8) plugins.RPC {
	a.method.Level = level
	return a
}

//PUT PUT路由
func (a *RPC) Method(method string, explain string, fn interface{}) {
	if a.method.Level <= 0 {
		a.method.Level = 1
	}
	if a.class != "" {
		method = a.class + "." + method
	}
	a.s.RegisterRPC(method, a.method.Level, a.method.IsAuth, explain, fn)
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
	//鉴权方法
	authMethod map[string]string
	//api信息
	api *serviceinfo.ServiceInfo
	//rpc服务信息
	rpc *serviceinfo.ServiceInfo
	//客户端
	client plugins.Client
	//参数编码器
	codec plugins.Codec
	//退出信号
	close int32
	//api注册拦截器
	apiRegIntercept func(gate, group, url string, level int8, isAuth bool, explain string)
	//meterics统计的数组
	metricValue []*metrics.MetricValue
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
		if service.cfg.GetDiscoveryModel() == config.NacosDiscoveryModel {
			service.register = nacosRegister.NewNacosRegister(service.cfg)
		}
		//使用consul
		if service.cfg.GetDiscoveryModel() == config.ConsulDiscoveryModel {
			service.register = consulRegister.NewConsulRegister(service.cfg)
		}
	}
	if service.router == nil {
		//默认路由
		service.router = router.NewRouter()
	}
	if service.limit == nil {
		//默认限流插件
		service.limit = limit.NewLimit(service.cfg.GetMaxServiceLimitRequest())
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
	//注册rpc服务
	service.rpc = &serviceinfo.ServiceInfo{
		Name:    service.name,
		Address: service.cfg.GetHost(),
		Port:    service.cfg.GetRPCPort(),
		Explain: service.cfg.GetExplain(),
		Time:    time.Now().Format("2006-01-02 15:04:05"),
	}
	//注册http服务
	service.api = &serviceinfo.ServiceInfo{
		Name:    service.name,
		Address: service.cfg.GetHost(),
		Port:    service.cfg.GetHTTPPort(),
		Explain: service.cfg.GetExplain(),
		Time:    time.Now().Format("2006-01-02 15:04:05"),
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

//RegisterRPC 注册RPC方法
func (s *Service) RegisterRPC(name string, level int8, isAuth bool, explain string, fn interface{}) {
	req, rep := s.router.RegisterByMethod(name, fn)
	method := &serviceinfo.Method{
		Name:     name,
		Level:    level,
		Explain:  explain,
		IsAuth:   isAuth,
		Request:  req,
		Response: rep,
	}
	s.rpc.Methods = append(s.rpc.Methods, method)
	if isAuth {
		s.authMethod[strings.ToLower(name)] = name
	}
	log.Traceln("注册RPC方法:", method.Name, "说明:", method.Explain)
}

//HTTP 创建http
func (s *Service) HTTP(gate string) plugins.HTTP {
	api := new(serviceinfo.API)
	api.Gate = gate
	return newHTTP(s, api)
}

//RPC 创建rpc
func (s *Service) RPC() plugins.RPC {
	method := new(serviceinfo.Method)
	return newRPC(s, method)
}

//APIRegIntercept API注册拦截器
func (s *Service) APIRegIntercept(f func(gate, group, url string, level int8, isAuth bool, explain string)) {
	s.apiRegIntercept = f
}

//AddMetricValue 添加metric采集的值
func (s *Service) AddMetricValue(metricValue []*metrics.MetricValue) {
	s.metricValue = append(s.metricValue, metricValue...)
}

//RegisterApi 注册API方法--注册给网管
func (s *Service) RegisterAPI(gate, group, methodname, version, path string, kind plugins.HTTPKind, level int8, isAuth bool, explain string, fn interface{}) {
	req, rep := s.router.RegisterByMethod(methodname, fn)
	url := fmt.Sprintf("/api/%s/%s/%s", s.name, version, path)
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
	s.api.API = append(s.api.API, api)
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
	//启动metrics
	if err := metrics.Init(&metrics.MetricOpts{
		NameSpace:     s.cfg.GetClusterName(),
		MetricsValues: s.metricValue,
	}); err != nil {
		log.Errorln(err.Error())
		return err
	}
	metrics.MetricServiceRun(s.name, 1)
	//监听指定信号
	c := make(chan os.Signal)
	defer close(c)
	signal.Notify(c, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		err := s.runTCP()
		if err != nil {
			log.Errorln(err.Error())
		}
		c <- nil
	}()
	go func() {
		err := s.runHTTP()
		if err != nil {
			log.Errorln(err.Error())
		}
		c <- nil
	}()
	log.Infoln("服务启动成功...")
	msg := <-c
	s.Close()
	metrics.MetricServiceRun(s.name, -1)
	if msg == nil {
		return fmt.Errorf("收到kill信号:%s", msg)
	}
	return nil
}

//_RunHTTP 启动HTTP
func (s *Service) runHTTP() error {
	//注册http接口服务
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(s.cors())
	router.Use(s.logger())
	if s.cfg.GetRunmode() == "trace" {
		pprof.Register(router)
	}
	router.GET("/apis", s.getAPI)
	router.GET("/rpc", s.getRPC)
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	httpport := fmt.Sprintf(":%d", s.cfg.GetHTTPPort())
	s.register.RegisterHTTPService(context.Background(), s.api)
	err := router.Run(httpport)
	if err != nil {
		log.Errorln(err.Error())
		return err
	}
	return nil
}

//getAPI 方法api信息
func (s *Service) getAPI(c *gin.Context) {
	c.JSON(http.StatusOK, s.api)
}

//getRPC 获取rpc服务信息
func (s *Service) getRPC(c *gin.Context) {
	c.JSON(http.StatusOK, s.rpc)
}

//logger 自定义日志输出
func (s *Service) logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()
		// 处理请求
		c.Next()
		// 结束时间
		end := time.Now()
		//执行时间
		latency := end.Sub(start)
		path := c.Request.URL.Path
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()
		log.Tracef("| %3d | %13v | %15s | %s  %s \n",
			statusCode,
			latency,
			clientIP,
			method,
			path,
		)
	}
}

//cors 处理跨域问题
func (s *Service) cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Headers", "Content-Type,TraceID, IsTest, Token,TimeOut")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type")
		c.Header("Access-Control-Allow-Credentials", "true")

		//放行所有OPTIONS方法
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		// 处理请求
		c.Next()
	}
}

//runTCP 启动TCP
func (s *Service) runTCP() error {
	l, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", s.cfg.GetRPCPort()))
	if err != nil {
		return err
	}
	defer l.Close()
	s.register.RegisterRPCService(context.Background(), s.rpc)
	for {
		if atomic.LoadInt32(&s.close) > 0 {
			return nil
		}
		conn, err := l.Accept()
		if err != nil {
			log.Traceln(err.Error())
			continue
		}
		go s.serveConn(conn)
	}
}

//log 日志
func (s *Service) log(address, name, method string, respone *header.Response) func() {
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
func (s *Service) serveConn(conn net.Conn) {
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
				defer s.log(req.Address, req.Name, req.Method, rep)()
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
