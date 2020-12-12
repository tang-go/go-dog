package mem

import (
	"sync"
	"time"
)

type mvalue struct {
	value string
	ttl   int64
}

//Mem 内存对象
type Mem struct {
	mem   sync.Map
	close chan bool
}

//NewMem 创建一个mem
func NewMem() *Mem {
	m := new(Mem)
	m.close = make(chan bool)
	go m.eventloop()
	return m
}

//eventloop 事件循环
func (m *Mem) eventloop() {
	for {
		select {
		case <-m.close:
			return
		case <-time.After(time.Millisecond):
			//一分钟清理一次内存
			m.mem.Range(func(key, v interface{}) bool {
				if value, o := v.(*mvalue); o {
					//超时了
					if value.ttl < time.Now().Unix() {
						m.mem.Delete(key)
					}
				} else {
					m.mem.Delete(key)
				}
				return true
			})
		}
	}
}

//Close 关闭
func (m *Mem) Close() {
	m.close <- true
}

//Set 设置内存缓存 ttl单位秒
func (m *Mem) Set(key string, value string, ttl int64) {
	m.mem.Store(key, &mvalue{value: value, ttl: time.Now().Unix() + ttl})
}

//Set 设置内存缓存
func (m *Mem) Get(key string) (string, bool) {
	v, ok := m.mem.Load(key)
	if ok {
		if value, o := v.(*mvalue); o {
			//超时了
			if value.ttl < time.Now().Unix() {
				m.mem.Delete(key)
				return "", false
			}
			return value.value, true
		}
	}
	return "", false
}

//Del 删除缓存
func (m *Mem) Del(key string) {
	m.mem.Delete(key)
}
