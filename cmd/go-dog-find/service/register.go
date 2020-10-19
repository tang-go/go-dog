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
	data        *param.Data
	service     *Service
}

//NewRegister 新建一个服务注册
func NewRegister(service *Service, conn net.Conn, data *param.Data, offlinefunc func()) *Register {
	return &Register{
		conn:        conn,
		offlinefunc: offlinefunc,
		data:        data,
		service:     service,
	}
}

//Run 启动
func (r *Register) Run() {
	defer r.offlinefunc()
	//上线服务
	r.service.cache.GetCache().Sadd(r.data.Label, r.data)
	//设置超时token
	r.service.cache.GetCache().SetByTime(r.data.Key, r.data, 5)
	for {
		_, _, err := io.ReadByTime(r.conn, time.Now().Add(time.Second*5))
		if err != nil {
			//删除服务
			r.service.cache.GetCache().SRem(r.data.Label, r.data)
			r.service.cache.GetCache().Del(r.data.Key)
			r.conn.Close()
			log.Errorln(err.Error())
			return
		}
		//收心跳
		r.service.cache.GetCache().SetByTime(r.data.Key, r.data, 5)
	}
}
