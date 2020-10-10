package gateway

import (
	"errors"
	"go-dog/cmd/define"
	customerror "go-dog/error"
	"go-dog/internal/client"
	"go-dog/internal/context"
	"go-dog/pkg/log"
	"go-dog/plugins"
	"go-dog/serviceinfo"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
)

type apiservice struct {
	method *serviceinfo.API
	name   string
	count  int32
}

//Gateway 服务发现
type Gateway struct {
	client   plugins.Client
	apis     map[string]*apiservice
	services map[string]*serviceinfo.APIServiceInfo
	lock     sync.RWMutex
}

//NewGateway  新建发现服务
func NewGateway() *Gateway {
	gateway := new(Gateway)
	gateway.apis = make(map[string]*apiservice)
	gateway.services = make(map[string]*serviceinfo.APIServiceInfo)
	//初始化客户端
	gateway.client = client.NewClient(define.TTL)
	//设置客户端最大访问量
	gateway.client.GetLimit().SetLimit(define.MaxClientRequestCount)
	//注册API上线事件
	gateway.client.GetDiscovery().RegAPIServiceOnlineNotice(gateway.APIServiceOnline)
	//注册API下线事件
	gateway.client.GetDiscovery().RegAPIServiceOfflineNotice(gateway.APIServiceOffline)
	//开启API事件监听
	gateway.client.GetDiscovery().WatchAPIService()
	return gateway
}

//APIServiceOnline api服务上线
func (g *Gateway) APIServiceOnline(key string, service *serviceinfo.APIServiceInfo) {
	g.lock.Lock()
	for _, method := range service.API {
		url := "/api/" + service.Name + "/" + method.Version + "/" + method.Path
		if api, ok := g.apis[url]; ok {
			api.count++
		} else {
			g.apis[url] = &apiservice{
				method: method,
				name:   service.Name,
				count:  1,
			}
			log.Traceln("收到API上线", key, method.Name, url, method.Request)
		}
		g.services[key] = service
	}
	g.lock.Unlock()
}

//APIServiceOffline api服务下线
func (g *Gateway) APIServiceOffline(key string) {
	g.lock.Lock()
	if service, ok := g.services[key]; ok {
		for _, method := range service.API {
			url := "/api/" + service.Name + "/" + method.Version + "/" + method.Path
			if api, ok := g.apis[url]; ok {
				api.count--
				if api.count <= 0 {
					delete(g.apis, url)
					log.Traceln("收到API下线", key, method.Name, url, method.Request)
				}
			}
		}
		delete(g.services, key)
	}
	g.lock.Unlock()
}

//Run 启动
func (g *Gateway) Run() {
	c := make(chan os.Signal)
	//监听指定信号
	signal.Notify(c, syscall.SIGINT, syscall.SIGKILL, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		gin.SetMode(gin.ReleaseMode)
		router := gin.New()
		router.Use(g.cors())
		router.Use(g.logger())
		//静态文件夹
		//router.StaticFS("/", http.Dir("./static"))
		//添加路由
		router.Any("/api/*router", g.routerResolution)
		log.Tracef("网管启动 0.0.0.0:80")
		err := router.Run(":80")
		if err != nil {
			panic(err.Error())
		}
	}()
	//阻塞直至有信号传入
	<-c
}

// routerResolution 路由解析
func (g *Gateway) routerResolution(c *gin.Context) {
	//路由解析
	url := c.Request.URL.String()
	g.lock.RLock()
	apiservice, ok := g.apis[url]
	g.lock.RUnlock()
	if !ok {
		c.JSON(http.StatusNotFound, customerror.EnCodeError(http.StatusNotFound, "路由错误"))
		return
	}
	if c.Request.Method != apiservice.method.Kind {
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
	body, err = g.validation(string(body), apiservice.method.Request)
	if err != nil {
		c.JSON(customerror.ParamError, customerror.EnCodeError(customerror.ParamError, err.Error()))
		return
	}
	ctx := context.Background()
	ctx.SetAddress(c.ClientIP())
	ctx.SetIsTest(isTest)
	ctx = context.WithTimeout(ctx, int64(time.Second*time.Duration(timeout)))
	back, err := g.client.SendRequest(ctx, plugins.RandomMode, apiservice.name, apiservice.method.Name, body)
	if err != nil {
		e := customerror.DeCodeError(err)
		c.JSON(e.Code, e)
		return
	}
	resp := new(interface{})
	g.client.GetCodec().DeCode(back, resp)
	c.JSON(http.StatusOK, gin.H{
		"Code": define.SuccessCode,
		"Body": resp,
		"Time": time.Now().Unix(),
	})
	return
}

//validation 验证参数
func (g *Gateway) validation(param string, tem map[string]interface{}) ([]byte, error) {
	p := make(map[string]interface{})
	if err := g.client.GetCodec().DeCode([]byte(param), &p); err != nil {
		return nil, err
	}
	if len(tem) != len(p) {
		return nil, errors.New("参数不正确")
	}
	for key := range p {
		if _, ok := tem[key]; !ok {
			return nil, errors.New("参数内容不正确")
		}
	}
	return g.client.GetCodec().EnCode(p)
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
		c.Header("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")
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
