package gateway

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	customerror "github.com/tang-go/go-dog/error"
	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/pkg/client"
	"github.com/tang-go/go-dog/pkg/config"
	"github.com/tang-go/go-dog/pkg/context"
	discovery "github.com/tang-go/go-dog/pkg/discovery/go-dog-find"
	"github.com/tang-go/go-dog/plugins"
)

//Gateway 服务发现
type Gateway struct {
	listenAPI             sync.Map
	client                plugins.Client
	cfg                   plugins.Cfg
	customGet             map[string]func(c *gin.Context)
	customPost            map[string]func(c *gin.Context)
	swaggerAuthCheck      func(token string) error
	authfunc              func(client plugins.Client, ctx plugins.Context, token, url string) error
	getRequestIntercept   func(c plugins.Context, url string, request []byte) ([]byte, bool, error)
	getResponseIntercept  func(c plugins.Context, url string, request []byte, response []byte)
	postRequestIntercept  func(c plugins.Context, url string, request []byte) ([]byte, bool, error)
	postResponseIntercept func(c plugins.Context, url string, request []byte, response []byte)
	discovery             *discovery.GoDogDiscovery
}

//NewGateway  新建发现服务
func NewGateway(name string) *Gateway {
	gateway := new(Gateway)
	//初始化配置
	gateway.cfg = config.NewConfig()
	//初始化服务发现
	gateway.discovery = discovery.NewGoDogDiscovery(gateway.cfg.GetDiscovery())
	gateway.discovery.WatchRPC()
	gateway.discovery.WatchAPI(name)
	gateway.discovery.ConnectClient()
	//初始化rpc服务
	gateway.client = client.NewClient(gateway.cfg, gateway.discovery)
	//初始化自定义请求
	gateway.customPost = make(map[string]func(c *gin.Context))
	gateway.customGet = make(map[string]func(c *gin.Context))
	//初始化文档
	return gateway
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
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(g.cors())
	router.Use(g.logger())
	for url, f := range g.customGet {
		router.GET(url, f)
	}
	for url, f := range g.customPost {
		router.POST(url, f)
	}
	//swagger 文档
	router.GET("/swagger/*any", g.getSwagger)
	//添加路由
	router.POST("/api/*router", g.routerPostResolution)
	//GET请求
	router.GET("/api/*router", g.routerGetResolution)
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
	token := c.Query("token")
	if c.Param("any") == "/swagger.json" {
		if token == "" {
			c.JSON(customerror.ParamError, customerror.EnCodeError(customerror.ParamError, "token不能为空"))
			return
		}
		if g.swaggerAuthCheck != nil {
			if err := g.swaggerAuthCheck(token); err != nil {
				c.JSON(customerror.ParamError, customerror.EnCodeError(customerror.ParamError, err.Error()))
				return
			}
		}
		c.String(200, g.ReadDoc())
		log.Traceln("获取swagger", token)
		return
	}
	ginSwagger.WrapHandler(swaggerFiles.Handler, func(c *ginSwagger.Config) {
		c.URL = "swagger.json?token=" + token
	})(c)
}

//routerGetResolution get路由解析
func (g *Gateway) routerGetResolution(c *gin.Context) {
	url := "/api" + c.Param("router")
	apiservice, ok := g.discovery.GetAPIByURL(url)
	if !ok {
		c.JSON(http.StatusNotFound, customerror.EnCodeError(http.StatusNotFound, "路由错误"))
		return
	}
	if c.Request.Method != apiservice.Method.Kind {
		c.JSON(http.StatusNotFound, customerror.EnCodeError(http.StatusNotFound, "路由错误"))
		return
	}

	timeoutstr := c.Request.Header.Get("TimeOut")
	if timeoutstr == "" {
		c.JSON(customerror.ParamError, customerror.EnCodeError(customerror.ParamError, "timeout不能为空"))
		return
	}
	timeout, err := strconv.Atoi(timeoutstr)
	if err != nil {
		c.JSON(customerror.ParamError, customerror.EnCodeError(customerror.ParamError, err.Error()))
		return
	}
	if timeout <= 0 {
		c.JSON(customerror.ParamError, customerror.EnCodeError(customerror.ParamError, "timeout必须大于0"))
		return
	}
	istest := c.Request.Header.Get("IsTest")
	if istest == "" {
		c.JSON(customerror.ParamError, customerror.EnCodeError(customerror.ParamError, "istest不能为空"))
		return
	}
	traceID := c.Request.Header.Get("TraceID")
	if traceID == "" {
		c.JSON(customerror.ParamError, customerror.EnCodeError(customerror.ParamError, "traceID不能为空"))
		return
	}
	isTest, err := strconv.ParseBool(istest)
	if err != nil {
		c.JSON(customerror.ParamError, customerror.EnCodeError(customerror.ParamError, err.Error()))
		return
	}

	p := make(map[string]interface{})
	for key, value := range apiservice.Method.Request {
		data := c.Query(key)
		if data == "" {
			c.JSON(customerror.ParamError, customerror.EnCodeError(customerror.ParamError, fmt.Sprintf("参数%s不正确", key)))
			return
		}
		vali, ok := value.(map[string]interface{})
		if !ok {
			c.JSON(customerror.ParamError, customerror.EnCodeError(customerror.ParamError, fmt.Sprintf("参数%v类型不正确", value)))
			return
		}
		tp, ok2 := vali["type"].(string)
		if !ok2 {
			c.JSON(customerror.ParamError, customerror.EnCodeError(customerror.ParamError, fmt.Sprintf("参数%v类型不是string", vali["type"])))
			return
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
	//查看方法是否需要验证权限
	if apiservice.Method.IsAuth {
		token := c.Request.Header.Get("Token")
		if token == "" {
			c.JSON(customerror.ParamError, customerror.EnCodeError(customerror.ParamError, "token不能为空"))
			return
		}
		//验证权限
		if g.authfunc != nil {
			if err := g.authfunc(g.GetClient(), ctx, token, url); err != nil {
				log.Errorln(err.Error())
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
				log.Errorln(err.Error())
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

// routerPostResolution post路由解析
func (g *Gateway) routerPostResolution(c *gin.Context) {
	//路由解析
	url := c.Request.URL.String()
	apiservice, ok := g.discovery.GetAPIByURL(url)
	if !ok {
		c.JSON(http.StatusNotFound, customerror.EnCodeError(http.StatusNotFound, "路由错误"))
		return
	}
	if c.Request.Method != apiservice.Method.Kind {
		c.JSON(http.StatusNotFound, customerror.EnCodeError(http.StatusNotFound, "路由错误"))
		return
	}
	timeoutstr := c.Request.Header.Get("TimeOut")
	if timeoutstr == "" {
		c.JSON(customerror.ParamError, customerror.EnCodeError(customerror.ParamError, "timeout不能为空"))
		return
	}
	timeout, err := strconv.Atoi(timeoutstr)
	if err != nil {
		c.JSON(customerror.ParamError, customerror.EnCodeError(customerror.ParamError, err.Error()))
		return
	}
	if timeout <= 0 {
		c.JSON(customerror.ParamError, customerror.EnCodeError(customerror.ParamError, "timeout必须大于0"))
		return
	}
	istest := c.Request.Header.Get("IsTest")
	if istest == "" {
		c.JSON(customerror.ParamError, customerror.EnCodeError(customerror.ParamError, "istest不能为空"))
		return
	}
	traceID := c.Request.Header.Get("TraceID")
	if traceID == "" {
		c.JSON(customerror.ParamError, customerror.EnCodeError(customerror.ParamError, "traceID不能为空"))
		return
	}
	isTest, err := strconv.ParseBool(istest)
	if err != nil {
		c.JSON(customerror.ParamError, customerror.EnCodeError(customerror.ParamError, err.Error()))
		return
	}

	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(customerror.ParamError, customerror.EnCodeError(customerror.ParamError, err.Error()))
		return
	}
	body, err = g.validation(string(body), apiservice.Method.Request)
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
	//查看方法是否需要验证权限
	if apiservice.Method.IsAuth {
		token := c.Request.Header.Get("Token")
		if token == "" {
			c.JSON(customerror.ParamError, customerror.EnCodeError(customerror.ParamError, "token不能为空"))
			return
		}
		//验证权限
		if g.authfunc != nil {
			if err := g.authfunc(g.GetClient(), ctx, token, url); err != nil {
				log.Errorln(err.Error())
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
				log.Errorln(err.Error())
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
func (g *Gateway) validation(param string, tem map[string]interface{}) ([]byte, error) {
	p := make(map[string]interface{})
	if err := g.GetClient().GetCodec().DeCode("json", []byte(param), &p); err != nil {
		return nil, err
	}
	for key := range p {
		if _, ok := tem[key]; !ok {
			log.Traceln("模版", tem, "传入参数", p)
			return nil, fmt.Errorf("不存在key为%s的参数", key)
		}
	}
	return g.GetClient().GetCodec().EnCode("json", p)
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
