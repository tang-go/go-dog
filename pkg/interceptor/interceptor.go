package interceptor

import (
	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/plugins"
	"time"
)

//Interceptor 拦截器
type Interceptor struct {
}

//NewInterceptor 创建一个拦截器
func NewInterceptor() *Interceptor {
	t := new(Interceptor)
	return t
}

//Request 请求
func (t *Interceptor) Request(ctx plugins.Context, servicename, method string, request interface{}) {
	ctx.SetData("Request-Time", time.Now())
}

//Respone 响应
func (t *Interceptor) Respone(ctx plugins.Context, servicename, method string, respone interface{}, err error) {
	start := ctx.GetDataByKey("Request-Time")
	tm, ok := start.(time.Time)
	if !ok {
		log.Errorln(start)
		return
	}
	end := time.Now()
	latency := end.Sub(tm)
	clientIP := ctx.GetAddress()
	log.Tracef("| %s | %13v | %s | %s ",
		clientIP,
		latency,
		servicename,
		method,
	)
}
