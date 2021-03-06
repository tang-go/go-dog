package fusing

import (
	"sync"
	"time"

	customerror "github.com/tang-go/go-dog/error"
	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/recover"
)

//方法错误统计
type method struct {
	name   string
	total  int64
	errnum int64
}

//Fusing 熔断模块
type Fusing struct {
	ttl     time.Duration
	methods map[string]*method
	auto    map[string]string
	forced  map[string]string
	err     map[string]error
	close   chan bool
	lock    sync.RWMutex
}

//NewFusing 新建一个熔断模块
func NewFusing(ttl time.Duration) *Fusing {
	fulsing := new(Fusing)
	fulsing.ttl = ttl
	fulsing.methods = make(map[string]*method)
	fulsing.forced = make(map[string]string)
	fulsing.auto = make(map[string]string)
	fulsing.close = make(chan bool)
	fulsing.err = make(map[string]error)
	go fulsing.eventloop()
	return fulsing
}

//SetFusingTTL 设置熔断统计时间
func (f *Fusing) SetFusingTTL(ttl time.Duration) {
	f.ttl = ttl
}

//AddError 添加服务错误
func (f *Fusing) AddError(servicekey string, err error) {
	f.lock.Lock()
	f.err[servicekey] = err
	f.lock.Unlock()
}

//AddErrorMethod 添加请求发生错误的方法
func (f *Fusing) AddErrorMethod(servicekey, methodname string, err error) {
	myError := customerror.DeCodeError(err)
	//只有系统错误才进入限流统计
	if myError.Code == customerror.RPCNotFind ||
		myError.Code == customerror.RequestTimeout ||
		myError.Code == customerror.InternalServerError ||
		myError.Code == customerror.ConnectClose ||
		myError.Code == customerror.SeviceLimitError {
		f.lock.Lock()
		if m, ok := f.methods[servicekey+"@"+methodname]; ok {
			m.errnum++
		}
		f.lock.Unlock()
		return
	}
}

//AddMethod 添加请求
func (f *Fusing) AddMethod(servicekey, methodname string) {
	f.lock.Lock()
	if m, ok := f.methods[servicekey+"@"+methodname]; ok {
		m.total++
	} else {
		f.methods[servicekey+"@"+methodname] = &method{
			name:   methodname,
			total:  1,
			errnum: 0,
		}
	}
	f.lock.Unlock()
}

//OpenFusing 设置某个服务方法强行开启熔断
func (f *Fusing) OpenFusing(servicekey, method string) {
	f.lock.Lock()
	f.forced[servicekey+"@"+method] = method
	log.Tracef("| 服务%s | 方法%s | 开启强制熔断 |", servicekey, method)
	f.lock.Unlock()
}

//CloseFusing 设置某个服务方法关闭熔断
func (f *Fusing) CloseFusing(servicekey, method string) {
	f.lock.Lock()
	delete(f.forced, servicekey+"@"+method)
	log.Tracef("| 服务%s | 方法%s | 关闭强制熔断 |", servicekey, method)
	f.lock.Unlock()
}

//IsFusing 是否熔断
func (f *Fusing) IsFusing(servicekey, method string) bool {
	f.lock.RLock()
	defer f.lock.RUnlock()
	if _, ok := f.err[servicekey]; ok {
		return true
	}
	if _, ok := f.forced[servicekey+"@"+method]; ok {
		return true
	}
	if _, ok := f.auto[servicekey+"@"+method]; ok {
		return true
	}
	return false
}

//Close 关闭
func (f *Fusing) Close() {
	f.close <- true
}

//eventloop 事件处理
func (f *Fusing) eventloop() {
	defer recover.Recover()
	for {
		select {
		case <-time.After(f.ttl):
			//清空所有统计数量
			f.lock.Lock()
			for key, m := range f.methods {
				if m.total > 10 && m.errnum > 0 {
					if m.errnum*100/m.total > 30 {
						f.auto[key] = m.name
						log.Tracef("| 服务%s | 方法%s | 开启自动熔断 |", key, m.name)
					} else {
						delete(f.auto, key)
					}
				} else {
					delete(f.auto, key)
				}
				m.errnum = 0
				m.total = 0
			}
			f.err = make(map[string]error)
			f.lock.Unlock()
		case <-f.close:
			close(f.close)
			return
		}
	}
}
