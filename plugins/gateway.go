package plugins

import "github.com/gin-gonic/gin"

//Gateway 网关
type Gateway interface {
	//SwaggerAuthCheck swagger权限检测
	SwaggerAuthCheck(swaggerAuthCheck func(token string) error)

	//GetRequestIntercept 拦截get请求
	GetRequestIntercept(f func(c Context, url string, request []byte) ([]byte, bool, error))

	//GetResponseIntercept 拦截get请求响应
	GetResponseIntercept(f func(c Context, url string, request []byte, response []byte))
	//PostRequestIntercept 拦截get请求
	PostRequestIntercept(f func(c Context, url string, request []byte) ([]byte, bool, error))

	//PostResponseIntercept 拦截get请求响应
	PostResponseIntercept(f func(c Context, url string, request []byte, response []byte))

	//OpenCustomGet 开启自定义get请求
	OpenCustomGet(url string, f func(c *gin.Context))

	//OpenCustomPost 开启自定义post请求
	OpenCustomPost(url string, f func(c *gin.Context))

	//GetClient 获取client
	GetClient() Client

	//GetCfg 获取cfg
	GetCfg() Cfg

	//Auth 验证权限
	Auth(f func(client Client, ctx Context, token, url string) error)

	//Run 启动
	Run(port int) error
}
