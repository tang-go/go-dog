package register

import (
	"context"
	"encoding/json"
	"fmt"
	"go-dog/cmd/go-dog-find/param"
	"go-dog/lib/io"
	"go-dog/log"
	"go-dog/serviceinfo"
	"net"
	"sync"
	"time"
)

//GoDogRegister 服务发现
type GoDogRegister struct {
	address    []string
	conn       net.Conn
	ttl        time.Duration
	pos        int
	count      int
	close      bool
	closeheart chan bool
	rpcinfo    *serviceinfo.ServiceInfo
	apiinfo    *serviceinfo.APIServiceInfo
	lock       sync.Mutex
}

//NewGoDogRegister  新建服务注册
func NewGoDogRegister(address []string, ttl int64) *GoDogRegister {
	dis := &GoDogRegister{
		address:    address,
		ttl:        2 * time.Second,
		count:      len(address),
		pos:        0,
		close:      false,
		closeheart: make(chan bool),
	}
	if err := dis._ConnectClient(); err != nil {
		panic(err)
	}
	return dis
}

//_ConnectClient 建立链接
func (d *GoDogRegister) _ConnectClient() error {
	d.lock.Lock()
	defer d.lock.Unlock()
	if d.close {
		return nil
	}
	address := d.address[d.pos]
	d.pos++
	if d.pos >= d.count {
		d.pos = 0
	}
	tcpAddr, err := net.ResolveTCPAddr("tcp4", address)
	if err != nil {
		return err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return err
	}
	//发送登陆请求
	login := new(param.LoginReq)
	login.Type = param.RegType
	buff, err := login.EnCode(login)
	if err != nil {
		conn.Close()
		log.Errorln(err.Error())
		return err
	}
	if _, err := io.WriteByTime(conn, buff, time.Now().Add(d.ttl)); err != nil {
		//断线开启重新链接
		conn.Close()
		log.Errorln(err.Error())
		return err
	}
	d.conn = conn
	go d._Watch()
	log.Traceln("链接成功注册中心", address)
	return nil
}

//_Watch 开始监听
func (d *GoDogRegister) _Watch() {
	go d._Heart()
	if d.rpcinfo != nil {
		d.RegisterRPCService(context.Background(), d.rpcinfo)
	}
	if d.apiinfo != nil {
		d.RegisterAPIService(context.Background(), d.apiinfo)
	}
	for {
		_, _, err := io.Read(d.conn)
		if err != nil {
			d.closeheart <- true
			d.conn.Close()
			log.Errorln(err.Error())
			break
		}
	}
	for {
		time.Sleep(d.ttl)
		log.Traceln("断线重链注册中心....")
		if d._ConnectClient() == nil {
			return
		}
	}
}

//RegisterRPCService 注册RPC服务
func (d *GoDogRegister) RegisterRPCService(ctx context.Context, info *serviceinfo.ServiceInfo) error {
	key := "rpc/" + fmt.Sprintf("%s:%d", info.Address, info.Port)
	info.Key = key
	val, err := json.Marshal(info)
	if err != nil {
		log.Errorln(err.Error())
		return err
	}
	event := &param.Event{
		Cmd:   param.Reg,
		Label: "/rpc",
		Data: &param.Data{
			Label: "/rpc",
			Key:   key,
			Value: string(val),
		},
	}
	buff, err := event.EnCode(event)
	if err != nil {
		log.Errorln(err.Error())
		return err
	}
	if _, err := io.WriteByTime(d.conn, buff, time.Now().Add(d.ttl)); err != nil {
		//断线开启重新链接
		d.conn.Close()
		log.Errorln(err.Error())
		return err
	}
	d.rpcinfo = info
	return nil
}

//RegisterAPIService 注册API服务
func (d *GoDogRegister) RegisterAPIService(ctx context.Context, info *serviceinfo.APIServiceInfo) error {
	key := "api/" + fmt.Sprintf("%s:%d", info.Address, info.Port)
	info.Key = key
	val, err := json.Marshal(info)
	if err != nil {
		log.Errorln(err.Error())
		return err
	}
	event := &param.Event{
		Cmd:   param.Reg,
		Label: "/api",
		Data: &param.Data{
			Label: "/api",
			Key:   key,
			Value: string(val),
		},
	}
	buff, err := event.EnCode(event)
	if err != nil {
		log.Errorln(err.Error())
		return err
	}
	if _, err := io.WriteByTime(d.conn, buff, time.Now().Add(d.ttl)); err != nil {
		//断线开启重新链接
		d.conn.Close()
		log.Errorln(err.Error())
		return err
	}
	d.apiinfo = info
	return nil
}

//_Heart 心跳
func (d *GoDogRegister) _Heart() {
	heart := &param.Event{
		Cmd: param.Heart,
	}
	buff, _ := heart.EnCode(heart)
	for {
		select {
		case <-d.closeheart:
			return
		case <-time.After(d.ttl):
			if _, err := io.WriteByTime(d.conn, buff, time.Now().Add(d.ttl)); err != nil {
				//断线开启重新链接
				d.conn.Close()
				log.Errorln(err.Error())
				break
			}
		}

	}
}

//Cancellation 注册服务
func (d *GoDogRegister) Cancellation() error {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.close = true
	if d.conn != nil {
		return d.conn.Close()
	}
	return nil
}
