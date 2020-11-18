package rpc

import (
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/tang-go/go-dog/header"
	"github.com/tang-go/go-dog/lib/io"
	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/plugins"
	"github.com/tang-go/go-dog/recover"
)

//ServiceRPC 服务
type ServiceRPC struct {
	codec      plugins.Codec
	conn       net.Conn
	isClose    int32
	callNotice func(*header.Request) *header.Response
	lock       sync.RWMutex
}

// NewServiceRPC 初始化一个service端rpc
func NewServiceRPC(conn net.Conn, codec plugins.Codec) *ServiceRPC {
	s := &ServiceRPC{
		conn:    conn,
		isClose: 0,
		codec:   codec,
	}
	go s.eventloop()
	return s
}

//RegisterCallNotice 注册client call通知
func (s *ServiceRPC) RegisterCallNotice(f func(*header.Request) *header.Response) {
	s.callNotice = f
}

// Close 关闭
func (s *ServiceRPC) Close() {
	atomic.AddInt32(&s.isClose, 1)
	s.conn.Close()
}

//Call 通知
func (s *ServiceRPC) call(req *header.Request) {
	defer recover.Recover()
	if s.callNotice != nil {
		rep := s.callNotice(req)
		s.send(rep)
	}
}

//Send 发送
func (s *ServiceRPC) send(response *header.Response) {
	if atomic.LoadInt32(&s.isClose) == 0 {
		buff, err := s.codec.EnCode("msgpack", response)
		if err != nil {
			return
		}
		_, err = io.Write(s.conn, buff)
		if err != nil {
			s.Close()
			return
		}
	}
}

//eventloop 事件监听
func (s *ServiceRPC) eventloop() {
	defer recover.Recover()
	defer func() {
		s.conn.Close()
	}()
	for {
		_, buff, err := io.ReadByTime(s.conn, time.Now().Add(time.Minute*5))
		if err != nil {
			s.Close()
			return
		}
		request := new(header.Request)
		err = s.codec.DeCode("msgpack", buff, request)
		if err != nil {
			log.Traceln(err.Error())
			continue
		}
		go s.call(request)
	}
}
