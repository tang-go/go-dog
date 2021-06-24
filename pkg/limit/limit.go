package limit

import (
	"time"

	"github.com/tang-go/go-dog/log"
	"golang.org/x/time/rate"
)

//Limit 限流
type Limit struct {
	max     int
	limiter *rate.Limiter
}

//NewLimit 创建一个默认限流插件
func NewLimit(max int) *Limit {
	limit := new(Limit)
	limit.limiter = rate.NewLimiter(rate.Every(time.Second/time.Duration(max)), max)
	log.Traceln("设置每秒最大流量", max)
	return limit
}

//SetLimit 设置最大限制
func (l *Limit) SetLimit(max int) {
	l.limiter.SetLimit(rate.Every(time.Second / time.Duration(max)))
	l.limiter.SetBurst(max)
	log.Traceln("设置每秒最大流量", max)
}

//IsLimit 获取是否可以通过
func (l *Limit) IsLimit() bool {
	return !l.limiter.Allow()
}

//Close 关闭
func (l *Limit) Close() {
}
