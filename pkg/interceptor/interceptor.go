package interceptor

import (
	"time"

	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/plugins"
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
	ctx.SetShare("Request-Time", time.Now())
}

//Respone 响应
func (t *Interceptor) Respone(ctx plugins.Context, servicename, method string, respone interface{}, err error) {
	start := ctx.GetShareByKey("Request-Time")
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
