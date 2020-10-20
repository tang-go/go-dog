package service

import (
	"go-dog/cmd/go-dog-find/param"
	"go-dog/lib/io"
	"go-dog/log"
	"net"
	"sync"
	"time"
)

//Discovery 服务发现
type Discovery struct {
	conn        net.Conn
	topics      map[string]string
	offlinefunc func()
	service     *Service
	close       chan bool
	lock        sync.RWMutex
}

//NewDiscovery 新建一个服务发现
func NewDiscovery(service *Service, conn net.Conn, offlinefunc func()) *Discovery {
	return &Discovery{
		conn:        conn,
		offlinefunc: offlinefunc,
		service:     service,
		close:       make(chan bool, 1),
		topics:      make(map[string]string),
	}
}

//Run 启动
func (d *Discovery) Run() {
	defer d.offlinefunc()
	go d.EventLoop()
	for {
		_, buff, err := io.ReadByTime(d.conn, time.Now().Add(time.Second*5))
		if err != nil {
			d.conn.Close()
			d.close <- true
			log.Errorln(err.Error())
			return
		}
		//处理各种事件
		event := new(param.Event)
		if err := event.DeCode(buff, event); err != nil {
			log.Errorln(err.Error())
			continue
		}
		switch event.Cmd {
		//监听事件
		case param.Listen:
			d.Listen(event.Label)
		}
	}
}

//Listen 监听
func (d *Discovery) Listen(label string) {
	d.lock.Lock()
	d.topics[label] = label
	log.Traceln("监听", label)
	d.lock.Unlock()
}

//PushEvent 推送事件
func (d *Discovery) PushEvent() {
	d.lock.RLock()
	for _, label := range d.topics {
		all := new(param.All)
		all.Label = label
		array, err := d.service.cache.GetCache().SMembers(label)
		if err == nil {
			for _, r := range array {
				data := new(param.Data)
				err := r.DeCode(data)
				if err != nil {
					log.Errorln(err.Error())
					continue
				}
				//查看服务是否还在线
				if err = d.service.cache.GetCache().Get(data.Key, data); err != nil {
					//不在线删除
					d.service.cache.GetCache().SRem(data.Label, data)
					log.Errorln(err.Error())
					continue
				}
				all.Datas = append(all.Datas, data)
			}
		}
		buff, err := all.EnCode(all)
		if err != nil {
			log.Errorln(err.Error())
			continue
		}
		//推送消息
		_, err = io.WriteByTime(d.conn, buff, time.Now().Add(time.Second*5))
		if err != nil {
			d.conn.Close()
			log.Errorln(err.Error())
			return
		}
	}
	d.lock.RUnlock()
}

//EventLoop 事件处理器
func (d *Discovery) EventLoop() {
	for {
		select {
		case <-d.close:
			close(d.close)
			return
		case <-time.After(time.Second * 1):
			//执行操作
			d.PushEvent()
		}
	}
}
