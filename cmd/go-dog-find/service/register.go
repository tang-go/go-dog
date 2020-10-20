package service

import (
	"go-dog/cmd/go-dog-find/param"
	"go-dog/lib/io"
	"go-dog/log"
	"net"
	"time"
)

//Register 注册
type Register struct {
	conn        net.Conn
	offlinefunc func()
	datas       []*param.Data
	service     *Service
}

//NewRegister 新建一个服务注册
func NewRegister(service *Service, conn net.Conn, offlinefunc func()) *Register {
	return &Register{
		conn:        conn,
		offlinefunc: offlinefunc,
		service:     service,
	}
}

//Run 启动
func (r *Register) Run() {
	defer r.offlinefunc()
	for {
		_, buff, err := io.ReadByTime(r.conn, time.Now().Add(time.Second*5))
		if err != nil {
			for _, data := range r.datas {
				r.service.cache.GetCache().SRem(data.Label, data)
				r.service.cache.GetCache().Del(data.Key)
			}
			r.conn.Close()
			log.Errorln(err.Error())
			return
		}
		event := new(param.Event)
		if err := event.DeCode(buff, event); err != nil {
			log.Errorln(err.Error())
			continue
		}
		if event.Cmd == param.Reg {
			//注册事件
			r.datas = append(r.datas, event.Data)
		}
		for _, data := range r.datas {
			//上线服务
			r.service.cache.GetCache().Sadd(data.Label, data)
			//设置超时token
			r.service.cache.GetCache().SetByTime(data.Key, data, 5)
		}
	}
}
