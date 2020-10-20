package rpc

import (
	customerror "github.com/tang-go/go-dog/error"
	"github.com/tang-go/go-dog/header"
	"github.com/tang-go/go-dog/lib/io"
	"github.com/tang-go/go-dog/lib/uuid"
	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/plugins"
	"github.com/tang-go/go-dog/recover"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type callmsg struct {
	response chan *header.Response
	resquest *header.Request
}

//ClientRPC 客户端
type ClientRPC struct {
	conn          net.Conn
	codec         plugins.Codec
	isClose       int32
	queue         map[string]*callmsg
	closecallback func(net.Conn)
	lock          sync.RWMutex
	wait          sync.WaitGroup
}

// NewClientRPC 创建一个msgpack客户端
func NewClientRPC(conn net.Conn, codec plugins.Codec, f func(net.Conn)) *ClientRPC {
	client := new(ClientRPC)
	client.conn = conn
	client.queue = make(map[string]*callmsg)
	client.closecallback = f
	client.isClose = 0
	client.codec = codec
	go client.eventloop()
	return client
}

//Call 调用函数
func (c *ClientRPC) Call(ctx plugins.Context, name, method string, request interface{}, response interface{}) (e error) {
	defer recover.Recover()
	c.wait.Add(1)
	defer c.wait.Done()
	if atomic.LoadInt32(&c.isClose) > 0 {
		return customerror.EnCodeError(customerror.ConnectClose, "链接已经关闭")
	}
	if ctx.GetTimeOut() < time.Now().UnixNano() {
		return customerror.EnCodeError(customerror.RequestTimeout, "请求超时")
	}
	req := new(header.Request)
	req.TTL = ctx.GetTTL()
	req.TimeOut = ctx.GetTimeOut()
	req.IsTest = ctx.GetIsTest()
	req.TraceID = ctx.GetTraceID()
	req.Address = ctx.GetAddress()
	req.Data = ctx.GetData()
	req.Token = ctx.GetToken()

	req.ID = uuid.GetToken()
	req.Name = name
	req.Method = method
	req.Arg, e = c.codec.EnCode("", request)
	if e != nil {
		return customerror.EnCodeError(customerror.ParamError, "参数不正确")
	}
	done := make(chan *header.Response, 1)
	c.wait.Add(1)
	go c.call(ctx, req, done)
	select {
	case rep := <-done:
		c.lock.Lock()
		delete(c.queue, req.ID)
		c.lock.Unlock()
		if rep.Error != nil {
			return rep.Error
		}
		c.codec.DeCode("", rep.Reply, response)
		return nil
	case <-ctx.Done():
		c.lock.Lock()
		delete(c.queue, req.ID)
		c.lock.Unlock()
		return customerror.EnCodeError(customerror.RequestTimeout, "请求超时")
	}
}

//SendRequest 发送请求
func (c *ClientRPC) SendRequest(ctx plugins.Context, name, method string, code string, arg []byte) (reply []byte, e error) {
	defer recover.Recover()
	c.wait.Add(1)
	defer c.wait.Done()
	if atomic.LoadInt32(&c.isClose) > 0 {
		return nil, customerror.EnCodeError(customerror.ConnectClose, "链接已经关闭")
	}
	if ctx.GetTimeOut() < time.Now().UnixNano() {
		return nil, customerror.EnCodeError(customerror.RequestTimeout, "请求超时")
	}
	req := new(header.Request)
	req.TTL = ctx.GetTTL()
	req.TimeOut = ctx.GetTimeOut()
	req.IsTest = ctx.GetIsTest()
	req.TraceID = ctx.GetTraceID()
	req.Address = ctx.GetAddress()
	req.Data = ctx.GetData()
	req.Token = ctx.GetToken()

	req.ID = uuid.GetToken()
	req.Name = name
	req.Method = method
	req.Arg = arg
	req.Code = code
	done := make(chan *header.Response, 1)
	c.wait.Add(1)
	go c.call(ctx, req, done)
	select {
	case rep := <-done:
		c.lock.Lock()
		delete(c.queue, req.ID)
		c.lock.Unlock()
		if rep.Error != nil {
			return nil, rep.Error
		}
		return rep.Reply, nil
	case <-ctx.Done():
		c.lock.Lock()
		delete(c.queue, req.ID)
		c.lock.Unlock()
		return nil, customerror.EnCodeError(customerror.RequestTimeout, "请求超时")
	}
}

//Call 调用函数
func (c *ClientRPC) call(ctx plugins.Context, req *header.Request, response chan *header.Response) {
	defer recover.Recover()
	defer c.wait.Done()
	if atomic.LoadInt32(&c.isClose) != 0 {
		rep := new(header.Response)
		rep.ID = req.ID
		rep.Name = req.Name
		rep.Error = customerror.EnCodeError(customerror.ConnectClose, "链接已经关闭")
		rep.Method = req.Method
		response <- rep
	}
	c.lock.RLock()
	if _, ok := c.queue[req.ID]; ok {
		rep := new(header.Response)
		rep.ID = req.ID
		rep.Name = req.Name
		rep.Error = customerror.EnCodeError(customerror.InternalServerError, "此请求ID已经存在,请勿重复请求")
		rep.Method = req.Method
		response <- rep
	}
	c.lock.RUnlock()
	c.send(req, response)
}

func (c *ClientRPC) done(response *header.Response) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	done, ok := c.queue[response.ID]
	if ok {
		done.response <- response
	}
}

func (c *ClientRPC) send(request *header.Request, response chan *header.Response) {
	if atomic.LoadInt32(&c.isClose) == 0 {
		c.lock.Lock()
		c.queue[request.ID] = &callmsg{
			resquest: request,
			response: response,
		}
		c.lock.Unlock()
		if request.TimeOut < time.Now().UnixNano() {
			//超时请求就没必要发起了
			return
		}
		buff, err := c.codec.EnCode("msgpack", request)
		if err == nil {
			_, err = io.Write(c.conn, buff)
			if err != nil {
				c.Close()
			}
		} else {
			rep := new(header.Response)
			rep.ID = request.ID
			rep.Name = request.Name
			rep.Error = customerror.EnCodeError(customerror.ParamError, "请求参数不正确")
			rep.Method = request.Method
			response <- rep
		}
	}
}

//eventloop 事件监听
func (c *ClientRPC) eventloop() {
	defer recover.Recover()
	defer func() {
		c.conn.Close()
		//给所有没有结束队列的请求返回失败
		c.lock.RLock()
		for _, vali := range c.queue {
			response := new(header.Response)
			response.ID = vali.resquest.ID
			response.Name = vali.resquest.Name
			response.Error = customerror.EnCodeError(customerror.InternalServerError, "服务链接已经关闭")
			response.Method = vali.resquest.Method
			vali.response <- response
		}
		c.lock.RUnlock()
		c.wait.Wait()
		if c.closecallback != nil {
			c.closecallback(c.conn)
		}
		log.Traceln("链接关闭", c.conn.RemoteAddr())
	}()
	for {
		_, buff, err := io.ReadByTime(c.conn, time.Now().Add(time.Minute*5))
		if err != nil {
			c.Close()
			return
		}
		response := new(header.Response)
		err = c.codec.DeCode("msgpack", buff, response)
		if err != nil {
			continue
		}
		c.done(response)
	}
}

//Close 关闭
func (c *ClientRPC) Close() {
	c.conn.Close()
	atomic.AddInt32(&c.isClose, 1)
}
