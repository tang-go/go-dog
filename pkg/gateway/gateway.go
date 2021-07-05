package gateway

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"text/template"
	"time"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	customerror "github.com/tang-go/go-dog/error"
	"github.com/tang-go/go-dog/jaeger"
	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/metrics"
	"github.com/tang-go/go-dog/pkg/client"
	"github.com/tang-go/go-dog/pkg/config"
	"github.com/tang-go/go-dog/pkg/context"
	consulRegister "github.com/tang-go/go-dog/pkg/discovery/consul"
	nacosDiscovery "github.com/tang-go/go-dog/pkg/discovery/nacos"
	"github.com/tang-go/go-dog/plugins"
	"github.com/tang-go/go-dog/serviceinfo"
)

//Gateway 服务发现
type Gateway struct {
	listenAPI             sync.Map
	name                  string
	client                plugins.Client
	cfg                   plugins.Cfg
	jaeger                *jaeger.Jaeger
	customGet             map[string]func(c *gin.Context)
	customPost            map[string]func(c *gin.Context)
	swaggerAuthCheck      func(token string) error
	authfunc              func(client plugins.Client, ctx plugins.Context, token, url string) error
	getRequestIntercept   func(c plugins.Context, url string, request []byte) ([]byte, bool, error)
	getResponseIntercept  func(c plugins.Context, url string, request []byte, response []byte)
	postRequestIntercept  func(c plugins.Context, url string, request []byte) ([]byte, bool, error)
	postResponseIntercept func(c plugins.Context, url string, request []byte, response []byte)
	discovery             plugins.Discovery
	metricValue           []*metrics.MetricValue
}

//NewGateway  新建发现服务
func NewGateway(name string) *Gateway {
	gateway := new(Gateway)
	gateway.name = name
	//初始化配置
	gateway.cfg = config.NewConfig()
	//初始化服务发现
	if gateway.cfg.GetModel() == config.NacosDiscoveryModel {
		gateway.discovery = nacosDiscovery.NewDiscovery(gateway.cfg)
	}
	if gateway.cfg.GetDiscoveryModel() == config.ConsulDiscoveryModel {
		gateway.discovery = consulRegister.NewDiscovery(gateway.cfg)
	}
	gateway.discovery.WatchAPI(name)
	//初始化rpc服务
	gateway.client = client.NewClient(gateway.cfg, gateway.discovery)
	//初始化自定义请求
	gateway.customPost = make(map[string]func(c *gin.Context))
	gateway.customGet = make(map[string]func(c *gin.Context))
	//初始化链路追踪
	gateway.jaeger = jaeger.NewJaeger(name, gateway.cfg)

	//默认监控
	labels := []string{Method, Path, Status}
	gateway.metricValue = []*metrics.MetricValue{
		{
			ValueType: metrics.Counter,
			Name:      ReqCount,
			Help:      "Counter. Total number of HTTP requests made",
			Labels:    labels,
		}, {
			ValueType: metrics.Histogram,
			Name:      ReqDuration,
			Help:      "Histogram. HTTP request latencies in seconds",
			Labels:    labels,
		}, {
			ValueType: metrics.Summary,
			Name:      ReqSizeBytes,
			Help:      "Summary. HTTP request sizes in bytes",
			Labels:    labels,
		}, {
			ValueType: metrics.Summary,
			Name:      RespSizeBytes,
			Help:      "Summary. HTTP request sizes in bytes",
			Labels:    labels,
		},
	}
	return gateway
}

//AddMetricValue 添加metric采集的值
func (g *Gateway) AddMetricValue(metricValue []*metrics.MetricValue) {
	g.metricValue = append(g.metricValue, metricValue...)
}

//SwaggerAuthCheck swagger权限检测
func (g *Gateway) SwaggerAuthCheck(swaggerAuthCheck func(token string) error) {
	g.swaggerAuthCheck = swaggerAuthCheck
}

//GetRequestIntercept 拦截get请求
func (g *Gateway) GetRequestIntercept(f func(c plugins.Context, url string, request []byte) ([]byte, bool, error)) {
	g.getRequestIntercept = f
}

//GetResponseIntercept 拦截get请求响应
func (g *Gateway) GetResponseIntercept(f func(c plugins.Context, url string, request []byte, response []byte)) {
	g.getResponseIntercept = f
}

//PostRequestIntercept 拦截get请求
func (g *Gateway) PostRequestIntercept(f func(c plugins.Context, url string, request []byte) ([]byte, bool, error)) {
	g.postRequestIntercept = f
}

//PostResponseIntercept 拦截get请求响应
func (g *Gateway) PostResponseIntercept(f func(c plugins.Context, url string, request []byte, response []byte)) {
	g.postResponseIntercept = f
}

//OpenCustomGet 开启自定义get请求
func (g *Gateway) OpenCustomGet(url string, f func(c *gin.Context)) {
	g.customGet[url] = f
}

//OpenCustomPost 开启自定义post请求
func (g *Gateway) OpenCustomPost(url string, f func(c *gin.Context)) {
	g.customPost[url] = f
}

//GetClient 获取client
func (g *Gateway) GetClient() plugins.Client {
	return g.client
}

//GetCfg 获取cfg
func (g *Gateway) GetCfg() plugins.Cfg {
	return g.cfg
}

//Auth 验证权限
func (g *Gateway) Auth(f func(client plugins.Client, ctx plugins.Context, token, url string) error) {
	g.authfunc = f
}

//Run 启动
func (g *Gateway) Run(port int) error {
	//启动metrics
	if err := metrics.Init(&metrics.MetricOpts{
		NameSpace:     g.cfg.GetClusterName(),
		SystemName:    strings.Replace(g.name, "/", "_", -1),
		MetricsValues: g.metricValue,
	}); err != nil {
		panic(err.Error())
	}
	//启动接口
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(g.cors())
	router.Use(g.logger())
	router.Use(metricMiddleware)
	for url, f := range g.customGet {
		router.GET(url, f)
	}
	for url, f := range g.customPost {
		router.POST(url, f)
	}
	if g.cfg.GetRunmode() == "trace" {
		pprof.Register(router)
	}
	router.GET("/metrics", gin.WrapH(promhttp.Handler()))
	router.GET("/swagger/*any", g.getSwagger)
	router.POST("/api/*router", g.routerPostAndPutResolution)
	router.PUT("/api/*router", g.routerPostAndPutResolution)
	router.GET("/api/*router", g.routerGetAndDeleteResolution)
	router.DELETE("/api/*router", g.routerGetAndDeleteResolution)
	c := make(chan os.Signal)
	//监听指定信号
	signal.Notify(c, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		httpport := fmt.Sprintf(":%d", port)
		log.Tracef("网管启动 0.0.0.0:%d", port)
		err := router.Run(httpport)
		if err != nil {
			panic(err.Error())
		}
	}()
	msg := <-c
	g.client.Close()
	return fmt.Errorf("收到kill信号:%s", msg)
}

//getSwagger 获取swagger
func (g *Gateway) getSwagger(c *gin.Context) {
	if c.Param("any") == "/swagger.json" {
		if g.swaggerAuthCheck != nil {
			token := c.Query("token")
			if token == "" {
				c.JSON(customerror.ParamError, customerror.EnCodeError(customerror.ParamError, "token不能为空"))
				return
			}
			if err := g.swaggerAuthCheck(token); err != nil {
				c.JSON(customerror.ParamError, customerror.EnCodeError(customerror.ParamError, err.Error()))
				return
			}
		}
		c.String(200, g.ReadDoc())
		return
	}
	ginSwagger.WrapHandler(swaggerFiles.Handler, func(c *ginSwagger.Config) {
		c.URL = "swagger.json"
	})(c)
}

//routerGetAndDeleteResolution get/delete路由解析
func (g *Gateway) routerGetAndDeleteResolution(c *gin.Context) {
	url := "/api" + c.Param("router")
	apiservice, ok := g.discovery.GetAPIByURL(c.Request.Method + url)
	if !ok {
		c.JSON(http.StatusNotFound, customerror.EnCodeError(http.StatusNotFound, "路由URL错误"))
		return
	}
	if c.Request.Method != apiservice.Method.Kind {
		c.JSON(http.StatusNotFound, customerror.EnCodeError(http.StatusNotFound, "路由Method错误"))
		return
	}
	timeoutstr := c.Request.Header.Get("timeOut")
	if timeoutstr == "" {
		timeoutstr = "6"
	}
	timeout, err := strconv.Atoi(timeoutstr)
	if err != nil {
		timeout = 6
	}
	if timeout <= 0 {
		c.JSON(customerror.ParamError, customerror.EnCodeError(customerror.ParamError, "timeout必须大于0"))
		return
	}
	istest := c.Request.Header.Get("isTest")
	if istest == "" {
		istest = "false"
	}
	isTest, err := strconv.ParseBool(istest)
	if err != nil {
		isTest = false
	}
	traceID := c.Request.Header.Get("traceID")
	if traceID == "" {
		c.JSON(customerror.ParamError, customerror.EnCodeError(customerror.ParamError, "traceID不能为空"))
		return
	}
	p := make(map[string]interface{})
	for key, value := range apiservice.Method.Request {
		vali, ok := value.(map[string]interface{})
		if !ok {
			c.JSON(customerror.ParamError, customerror.EnCodeError(customerror.ParamError, fmt.Sprintf("传入参数%v类型错误", value)))
			return
		}
		tp, ok2 := vali["type"].(string)
		if !ok2 {
			c.JSON(customerror.ParamError, customerror.EnCodeError(customerror.ParamError, fmt.Sprintf("传入参数%v类型错误", value)))
			return
		}
		required := false
		re, ok3 := vali["required"]
		if ok3 && re == "true" {
			required = true
		}
		data := c.Query(key)
		if data == "" {
			if required {
				c.JSON(customerror.ParamError, customerror.EnCodeError(customerror.ParamError, fmt.Sprintf("参数%s必须传入", key)))
				return
			}
			continue
		}
		v, err := transformation(tp, data)
		if err != nil {
			c.JSON(customerror.ParamError, customerror.EnCodeError(customerror.ParamError, err.Error()))
			return
		}
		p[key] = v
	}
	body, err := g.GetClient().GetCodec().EnCode("json", p)
	if err != nil {
		c.JSON(customerror.ParamError, customerror.EnCodeError(customerror.ParamError, err.Error()))
		return
	}

	ctx := context.Background()
	ctx.SetAddress(c.ClientIP())
	ctx.SetIsTest(isTest)
	ctx.SetTraceID(traceID)
	ctx.SetURL(url)
	ctx = context.WithTimeout(ctx, int64(time.Second*time.Duration(timeout)))
	//开启追踪
	if span, err := g.jaeger.StartSpan(ctx, url); err == nil {
		defer span.Finish()
	}
	//查看方法是否需要验证权限
	if apiservice.Method.IsAuth {
		token := c.Request.Header.Get("token")
		if token == "" {
			c.JSON(customerror.ParamError, customerror.EnCodeError(customerror.ParamError, "token不能为空"))
			return
		}
		//验证权限
		if g.authfunc != nil {
			if err := g.authfunc(g.GetClient(), ctx, token, url); err != nil {
				log.Traceln(err.Error())
				c.JSON(customerror.ParamError, customerror.EnCodeError(customerror.ParamError, "token不正确"))
				return
			}
		}
		//设置token
		ctx.SetToken(token)
	}
	//拦截请求
	if g.getRequestIntercept != nil {
		if reposne, ok, err := g.getRequestIntercept(ctx, url, body); ok {
			if err != nil {
				log.Traceln(err.Error())
				c.JSON(customerror.ParamError, customerror.EnCodeError(customerror.ParamError, err.Error()))
				return
			}
			resp := make(map[string]interface{})
			g.GetClient().GetCodec().DeCode("json", reposne, &resp)
			c.JSON(http.StatusOK, gin.H{
				"code": 10000,
				"body": resp,
				"time": time.Now().Unix(),
			})
			return
		}
	}
	back, err := g.GetClient().SendRequest(ctx, plugins.RandomMode, apiservice.Name, apiservice.Method.Name, "json", body)
	if err != nil {
		e := customerror.DeCodeError(err)
		c.JSON(http.StatusOK, e)
		return
	}
	//拦截返回
	if g.getResponseIntercept != nil {
		g.getResponseIntercept(ctx, url, body, back)
	}
	resp := make(map[string]interface{})
	g.GetClient().GetCodec().DeCode("json", back, &resp)
	c.JSON(http.StatusOK, gin.H{
		"code": 10000,
		"body": resp,
		"time": time.Now().Unix(),
	})
	return
}

// routerPostAndPutResolution post/put路由解析
func (g *Gateway) routerPostAndPutResolution(c *gin.Context) {
	//路由解析
	url := c.Request.URL.String()
	apiservice, ok := g.discovery.GetAPIByURL(c.Request.Method + url)
	if !ok {
		c.JSON(http.StatusNotFound, customerror.EnCodeError(http.StatusNotFound, "路由URL错误"))
		return
	}
	if c.Request.Method != apiservice.Method.Kind {
		c.JSON(http.StatusNotFound, customerror.EnCodeError(http.StatusNotFound, "路由Method错误"))
		return
	}
	timeoutstr := c.Request.Header.Get("timeOut")
	if timeoutstr == "" {
		timeoutstr = "6"
	}
	timeout, err := strconv.Atoi(timeoutstr)
	if err != nil {
		timeout = 6
	}
	if timeout <= 0 {
		c.JSON(customerror.ParamError, customerror.EnCodeError(customerror.ParamError, "timeout必须大于0"))
		return
	}
	istest := c.Request.Header.Get("isTest")
	if istest == "" {
		istest = "false"
	}
	isTest, err := strconv.ParseBool(istest)
	if err != nil {
		isTest = false
	}
	traceID := c.Request.Header.Get("traceID")
	if traceID == "" {
		c.JSON(customerror.ParamError, customerror.EnCodeError(customerror.ParamError, "traceID不能为空"))
		return
	}

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(customerror.ParamError, customerror.EnCodeError(customerror.ParamError, err.Error()))
		return
	}
	if err := g.validation(body, apiservice.Method.Request); err != nil {
		c.JSON(customerror.ParamError, customerror.EnCodeError(customerror.ParamError, err.Error()))
		return
	}
	ctx := context.Background()
	ctx.SetAddress(c.ClientIP())
	ctx.SetIsTest(isTest)
	ctx.SetTraceID(traceID)
	ctx.SetURL(url)
	ctx = context.WithTimeout(ctx, int64(time.Second*time.Duration(timeout)))
	//开启追踪
	if span, err := g.jaeger.StartSpan(ctx, url); err == nil {
		defer span.Finish()
	}
	//查看方法是否需要验证权限
	if apiservice.Method.IsAuth {
		token := c.Request.Header.Get("token")
		if token == "" {
			c.JSON(customerror.ParamError, customerror.EnCodeError(customerror.ParamError, "token不能为空"))
			return
		}
		//验证权限
		if g.authfunc != nil {
			if err := g.authfunc(g.GetClient(), ctx, token, url); err != nil {
				log.Traceln(err.Error())
				c.JSON(customerror.ParamError, customerror.EnCodeError(customerror.ParamError, "token不正确"))
				return
			}
		}
		//设置token
		ctx.SetToken(token)
	}
	//拦截请求
	if g.postRequestIntercept != nil {
		if reposne, ok, err := g.postRequestIntercept(ctx, url, body); ok {
			if err != nil {
				log.Traceln(err.Error())
				c.JSON(customerror.ParamError, customerror.EnCodeError(customerror.ParamError, err.Error()))
				return
			}
			resp := make(map[string]interface{})
			g.GetClient().GetCodec().DeCode("json", reposne, &resp)
			c.JSON(http.StatusOK, gin.H{
				"code": 10000,
				"body": resp,
				"time": time.Now().Unix(),
			})
			return
		}
	}
	back, err := g.GetClient().SendRequest(ctx, plugins.RandomMode, apiservice.Name, apiservice.Method.Name, "json", body)
	if err != nil {
		e := customerror.DeCodeError(err)
		c.JSON(http.StatusOK, e)
		return
	}
	//拦截返回
	if g.postResponseIntercept != nil {
		g.postResponseIntercept(ctx, url, body, back)
	}
	resp := make(map[string]interface{})
	g.GetClient().GetCodec().DeCode("json", back, &resp)
	c.JSON(http.StatusOK, gin.H{
		"code": 10000,
		"body": resp,
		"time": time.Now().Unix(),
	})
	return
}

//validation 验证参数
func (g *Gateway) validation(param []byte, tem map[string]interface{}) error {
	p := make(map[string]interface{})
	if err := g.GetClient().GetCodec().DeCode("json", param, &p); err != nil {
		return err
	}
	for key := range p {
		if _, ok := tem[key]; !ok {
			log.Traceln("模版", tem, "传入参数", p)
			return fmt.Errorf("不存在值为%s的参数", key)
		}
	}
	return nil
}

//logger 自定义日志输出
func (g *Gateway) logger() gin.HandlerFunc {
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
		log.Infof("| %3d | %13v | %15s | %s  %s \n",
			statusCode,
			latency,
			clientIP,
			method,
			path,
		)
	}
}

//cors 处理跨域问题
func (g *Gateway) cors() gin.HandlerFunc {
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

//Docs 文档内容
type Docs struct {
	Swagger     string                            `json:"swagger"`
	Info        Info                              `json:"info"`
	Host        string                            `json:"host"`
	BasePath    string                            `json:"basePath"`
	Paths       map[string]map[string]interface{} `json:"paths"`
	Definitions map[string]Definitions            `json:"definitions"`
}

//Info 信息
type Info struct {
	Description string `json:"description"`
	Title       string `json:"title"`
	Contact     struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"contact"`
	License struct {
	} `json:"license"`
	Version string `json:"version"`
}

//Definitions 参数定义
type Definitions struct {
	Name       string                 `json:"name"`
	Type       string                 `json:"type"`
	Properties map[string]Description `json:"properties"`
}

//Description 描述
type Description struct {
	Description string      `json:"description"`
	Type        string      `json:"type"`
	Items       interface{} `json:"items,omitempty"`
	Ref         string      `json:"$ref,omitempty"`
}

//Ref 链接
type Ref struct {
	Ref string `json:"$ref,omitempty"`
}

//Body 请求
type Body struct {
	Consumes   []string     `json:"consumes"`
	Produces   []string     `json:"produces"`
	Tags       []string     `json:"tags"`
	Summary    string       `json:"summary"`
	Parameters []Parameters `json:"parameters"`
	Responses  struct {
		Code200 struct {
			Description string `json:"description"`
			Schema      struct {
				Type string `json:"type"`
				Ref  string `json:"$ref,omitempty"`
			} `json:"schema"`
		} `json:"200"`
	} `json:"responses"`
}

//Parameters api描述
type Parameters struct {
	Type        string `json:"type,omitempty"`
	Description string `json:"description"`
	Name        string `json:"name"`
	In          string `json:"in"`
	Required    bool   `json:"required"`
	Schema      struct {
		Type string `json:"type"`
		Ref  string `json:"$ref,omitempty"`
	} `json:"schema"`
}

//t type解析
func t(tp string) string {
	switch tp {
	case "int8":
		return "integer"
	case "int":
		return "integer"
	case "int32":
		return "integer"
	case "int64":
		return "integer"
	case "uint8":
		return "integer"
	case "uint":
		return "integer"
	case "uint32":
		return "integer"
	case "uint64":
		return "integer"
	case "float":
		return "number"
	case "float32":
		return "number"
	case "float64":
		return "number"
	case "byte":
		return "string"
	case "bool":
		return "boolean"
	default:
		return tp
	}
}

//transformation 转换
func transformation(tp string, value string) (interface{}, error) {
	switch tp {
	case "int8":
		i, e := strconv.ParseInt(value, 10, 8)
		if e != nil {
			return nil, fmt.Errorf("需要参数是int8 %s是", e.Error())
		}
		return int8(i), nil
	case "int":
		return strconv.Atoi(value)
	case "int32":
		i, e := strconv.ParseInt(value, 10, 32)
		if e != nil {
			return nil, fmt.Errorf("需要参数是int32 %s是", e.Error())
		}
		return int32(i), nil
	case "int64":
		return strconv.ParseInt(value, 10, 64)
	case "uint8":
		i, e := strconv.ParseInt(value, 10, 8)
		if e != nil {
			return nil, e
		}
		return uint8(i), nil
	case "uint":
		i, e := strconv.Atoi(value)
		if e != nil {
			return nil, fmt.Errorf("需要参数是uint %s是", e.Error())
		}
		return uint(i), nil
	case "uint32":
		i, e := strconv.ParseInt(value, 10, 32)
		if e != nil {
			return nil, fmt.Errorf("需要参数是uint32 %s是", e.Error())
		}
		return uint32(i), nil
	case "uint64":
		i, e := strconv.ParseInt(value, 10, 64)
		if e != nil {
			return nil, fmt.Errorf("需要参数是uint64 %s是", e.Error())
		}
		return uint64(i), nil
	case "float32":
		return strconv.ParseFloat(value, 32)
	case "float64":
		return strconv.ParseFloat(value, 64)
	case "bool":
		return strconv.ParseBool(value)
	case "string":
		return value, nil
	default:
		return tp, fmt.Errorf("暂时不支持此类型参数%s", tp)
	}
}

//createPostAndPutAPI 创建一个POST/PUT API
func createPostAndPutAPI(kind, tags, summary, name string, isAuth bool, request, respone map[string]interface{}) (a Body, definitions []Definitions) {
	api := Body{
		Consumes: []string{"application/json"},
		Produces: []string{"application/json"},
		Tags:     []string{tags},
		Summary:  summary,
	}
	parameters := Parameters{
		Description: "请求内容",
		Name:        "body",
		In:          "body",
		Required:    true,
	}
	requestName := strings.Replace(tags+"."+name+"."+kind+".Request", "/", ".", -1)
	requestProperties := createDefinitions(requestName, request)
	definitions = append(definitions, requestProperties...)

	parameters.Schema.Type = "object"
	parameters.Schema.Ref = "#/definitions/" + requestName
	api.Parameters = []Parameters{
		{
			Type:        "integer",
			Description: "请求超时时间,单位秒",
			Name:        "timeOut",
			In:          "header",
			Required:    true,
		},
		{
			Type:        "string",
			Description: "链路请求ID",
			Name:        "traceID",
			In:          "header",
			Required:    true,
		},
		{
			Type:        "boolean",
			Description: "是否是测试请求",
			Name:        "isTest",
			In:          "header",
			Required:    true,
		},
	}
	if isAuth {
		api.Parameters = append(api.Parameters, Parameters{
			Type:        "string",
			Description: "验证Token",
			Name:        "token",
			In:          "header",
			Required:    true,
		})
	}
	api.Parameters = append(api.Parameters, parameters)

	responeName := strings.Replace(tags+"."+name+"."+kind+".Respone", "/", ".", -1)
	responeProperties := createDefinitions(responeName, respone)
	definitions = append(definitions, responeProperties...)

	api.Responses.Code200.Description = "请求成功返回参数"
	api.Responses.Code200.Schema.Type = "object"
	api.Responses.Code200.Schema.Ref = "#/definitions/" + responeName

	return api, definitions
}

//createGetAndDeleteAPI 创建一个GET/DELETE API
func createGetAndDeleteAPI(kind, tags, summary, name string, isAuth bool, request, respone map[string]interface{}) (a Body, definitions []Definitions) {
	api := Body{
		Consumes: []string{"application/json"},
		Tags:     []string{tags},
		Summary:  summary,
	}
	for key, value := range request {
		if vali, ok := value.(map[string]interface{}); ok {
			des, ok1 := vali["description"]
			tp, ok2 := vali["type"]
			required := false
			re, ok3 := vali["required"]
			if ok3 && re == "true" {
				required = true
			}
			if ok1 == true && ok2 == true {
				api.Parameters = append(api.Parameters, Parameters{
					Type:        t(tp.(string)),
					Description: des.(string),
					Name:        key,
					In:          "query",
					Required:    required,
				})
			}
		}
	}
	api.Parameters = append(api.Parameters,
		Parameters{
			Type:        "integer",
			Description: "请求超时时间,单位秒",
			Name:        "timeOut",
			In:          "header",
			Required:    true,
		},
		Parameters{
			Type:        "string",
			Description: "链路请求ID",
			Name:        "traceID",
			In:          "header",
			Required:    true,
		},
		Parameters{
			Type:        "boolean",
			Description: "是否是测试请求",
			Name:        "isTest",
			In:          "header",
			Required:    true,
		})
	if isAuth {
		api.Parameters = append(api.Parameters, Parameters{
			Type:        "string",
			Description: "验证Token",
			Name:        "token",
			In:          "header",
			Required:    true,
		})
	}

	responeName := strings.Replace(tags+"."+name+"."+kind+".Respone", "/", ".", -1)
	responeProperties := createDefinitions(responeName, respone)
	definitions = append(definitions, responeProperties...)

	api.Responses.Code200.Description = "请求成功返回参数"
	api.Responses.Code200.Schema.Type = "object"
	api.Responses.Code200.Schema.Ref = "#/definitions/" + responeName

	return api, definitions
}

//createDefinitions 生成Definitions
func createDefinitions(name string, mp map[string]interface{}) (definitions []Definitions) {
	properties := make(map[string]Description)
	for key, value := range mp {
		if vali, ok := value.(map[string]interface{}); ok {
			slice, ok := vali["slice"]
			des, ok1 := vali["description"]
			tp, ok2 := vali["type"]
			if ok {
				mp, o := slice.(map[string]interface{})
				if o == true {
					description := Description{}
					if ok1 {
						description.Description = des.(string)
					}
					if ok2 {
						description.Type = t(tp.(string))
					}
					son := name + "." + key
					definitions = append(definitions, createDefinitions(son, mp)...)
					description.Items = &Ref{
						Ref: "#/definitions/" + son,
					}
					properties[key] = description
				} else {
					description := Description{}
					if ok1 {
						description.Description = des.(string)
					}
					if ok2 {
						description.Type = t(tp.(string))
					}
					description.Items = map[string]string{
						"type": t(vali["slice"].(string)),
					}
					properties[key] = description
				}
				continue
			} else if object, ok3 := vali["object"]; ok3 {
				mp, o := object.(map[string]interface{})
				if o == true {
					description := Description{}
					if ok1 {
						description.Description = des.(string)
					}
					description.Type = "object"
					son := name + "." + key
					definitions = append(definitions, createDefinitions(son, mp)...)
					description.Ref = "#/definitions/" + son

					properties[key] = description
					continue
				}
			}
			description := Description{}
			if ok1 {
				description.Description = des.(string)
			}
			if ok2 {
				description.Type = t(tp.(string))
			}
			properties[key] = description
		}
	}
	definition := Definitions{
		Name:       name,
		Type:       "object",
		Properties: properties,
	}
	definitions = append(definitions, definition)
	return
}

//swagger info
type swaggerInfo struct {
	Version     string
	Host        string
	BasePath    string
	Schemes     []string
	Title       string
	Description string
}

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = swaggerInfo{Schemes: []string{}}

//assembleDocs 组装文档
func (g *Gateway) assembleDocs() string {
	info := Info{
		Description: "",
		Title:       g.name + "网管API文档",
		Version:     "{{.Version}}",
	}

	paths := make(map[string]map[string]interface{})
	definitions := make(map[string]Definitions)

	g.discovery.RangeAPI(func(url string, service *serviceinfo.ServcieAPI) {
		switch service.Method.Kind {
		case string(plugins.POST):
			api, d := createPostAndPutAPI(
				service.Method.Kind,
				service.Explain+"["+service.Tags+"]",
				service.Method.Explain,
				service.Method.Name,
				service.Method.IsAuth,
				service.Method.Request,
				service.Method.Response)
			if value, ok := paths[service.Method.Path]; ok {
				value["post"] = api
				paths[service.Method.Path] = value
			} else {
				value := make(map[string]interface{})
				value["post"] = api
				paths[service.Method.Path] = value
			}
			for _, definition := range d {
				definitions[definition.Name] = definition
			}
		case string(plugins.PUT):
			api, d := createPostAndPutAPI(
				service.Method.Kind,
				service.Explain+"["+service.Tags+"]",
				service.Method.Explain,
				service.Method.Name,
				service.Method.IsAuth,
				service.Method.Request,
				service.Method.Response)
			if value, ok := paths[service.Method.Path]; ok {
				value["put"] = api
				paths[service.Method.Path] = value
			} else {
				value := make(map[string]interface{})
				value["put"] = api
				paths[service.Method.Path] = value
			}
			for _, definition := range d {
				definitions[definition.Name] = definition
			}
		case string(plugins.GET):
			api, d := createGetAndDeleteAPI(
				service.Method.Kind,
				service.Explain+"["+service.Tags+"]",
				service.Method.Explain,
				service.Method.Name,
				service.Method.IsAuth,
				service.Method.Request,
				service.Method.Response)
			if value, ok := paths[service.Method.Path]; ok {
				value["get"] = api
				paths[service.Method.Path] = value
			} else {
				value := make(map[string]interface{})
				value["get"] = api
				paths[service.Method.Path] = value
			}
			for _, definition := range d {
				definitions[definition.Name] = definition
			}
		case string(plugins.DELETE):
			api, d := createGetAndDeleteAPI(
				service.Method.Kind,
				service.Explain+"["+service.Tags+"]",
				service.Method.Explain,
				service.Method.Name,
				service.Method.IsAuth,
				service.Method.Request,
				service.Method.Response)
			if value, ok := paths[service.Method.Path]; ok {
				value["delete"] = api
				paths[service.Method.Path] = value
			} else {
				value := make(map[string]interface{})
				value["delete"] = api
				paths[service.Method.Path] = value
			}
			for _, definition := range d {
				definitions[definition.Name] = definition
			}
		}
	})

	docs := &Docs{
		Swagger:     "2.0",
		Host:        "{{.Host}}",
		BasePath:    "{{.BasePath}}",
		Info:        info,
		Paths:       paths,
		Definitions: definitions,
	}
	buff, _ := json.Marshal(docs)
	return string(buff)
}

//ReadDoc 读取文档
func (g *Gateway) ReadDoc() string {
	docs := g.assembleDocs()
	t, err := template.New("swagger_info").Funcs(template.FuncMap{
		"marshal": func(v interface{}) string {
			a, _ := json.Marshal(v)
			return string(a)
		},
	}).Parse(docs)
	if err != nil {
		log.Traceln(err.Error())
		return docs
	}
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, SwaggerInfo); err != nil {
		log.Traceln(err.Error())
		return docs
	}
	return tpl.String()
}
