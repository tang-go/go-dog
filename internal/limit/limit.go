package limit

import (
	"go-dog/recover"
	"sync/atomic"
	"time"
)

//Limit 限流
type Limit struct {
	max   int64
	count int64
	close chan bool
}

//NewLimit 创建一个默认限流插件
func NewLimit(max int64) *Limit {
	limit := new(Limit)
	limit.count = max
	limit.max = max
	limit.close = make(chan bool)
	go limit.eventloop()
	return limit
}

//SetLimit 设置最大限制
func (l *Limit) SetLimit(max int64) {
	atomic.StoreInt64(&l.max, max)
}

//IsLimit 获取是否可以通过
func (l *Limit) IsLimit() bool {
	atomic.AddInt64(&l.count, -1)
	if atomic.LoadInt64(&l.count) >= 0 {
		return false
	}
	return true
}

//Close 关闭
func (l *Limit) Close() {
	l.close <- true
}

//事件循环
func (l *Limit) eventloop() {
	defer recover.Recover()
	for {
		select {
		case <-time.After(time.Second * 1):
			atomic.StoreInt64(&l.count, atomic.LoadInt64(&l.max))
		case <-l.close:
			close(l.close)
			return
		}
	}
}
