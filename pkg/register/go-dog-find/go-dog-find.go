package register

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/tang-go/go-dog/lib/io"
	"github.com/tang-go/go-dog/lib/rand"
	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/pkg/register/go-dog-find/param"
	"github.com/tang-go/go-dog/serviceinfo"
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
	rpcinfo    *serviceinfo.RPCServiceInfo
	apiinfo    *serviceinfo.APIServiceInfo
	lock       sync.Mutex
}

//NewGoDogRegister  新建服务注册
func NewGoDogRegister(address []string) *GoDogRegister {
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
	index := rand.IntRand(0, d.count)
	address := d.address[index]
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
	if err := d._SendMsg(conn, param.Listen, buff); err != nil {
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

//_SendMsg 发送消息
func (d *GoDogRegister) _SendMsg(conn net.Conn, cmd int8, buff []byte) error {
	event := new(param.Event)
	event.Cmd = param.Login
	event.Data = buff
	if _, err := io.WriteByTime(conn, buff, time.Now().Add(d.ttl)); err != nil {
		log.Errorln(err.Error())
		return err
	}
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
func (d *GoDogRegister) RegisterRPCService(ctx context.Context, info *serviceinfo.RPCServiceInfo) error {
	key := "rpc/" + fmt.Sprintf("%s:%d", info.Address, info.Port)
	info.Key = key
	val, err := json.Marshal(info)
	if err != nil {
		log.Errorln(err.Error())
		return err
	}
	reg := &param.RegReq{
		Label: param.RPCLabel,
		Data: param.Data{
			Key:   key,
			Value: string(val),
		},
	}
	buff, err := reg.EnCode(reg)
	if err != nil {
		log.Errorln(err.Error())
		return err
	}
	if err := d._SendMsg(d.conn, param.Reg, buff); err != nil {
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
	reg := &param.RegReq{
		Label: param.APILabel,
		Data: param.Data{
			Key:   key,
			Value: string(val),
		},
	}
	buff, err := reg.EnCode(reg)
	if err != nil {
		log.Errorln(err.Error())
		return err
	}
	if err := d._SendMsg(d.conn, param.Reg, buff); err != nil {
		log.Errorln(err.Error())
		return err
	}
	d.apiinfo = info
	return nil
}

//_Heart 心跳
func (d *GoDogRegister) _Heart() {
	for {
		select {
		case <-d.closeheart:
			return
		case <-time.After(d.ttl):
			if err := d._SendMsg(d.conn, param.Heart, nil); err != nil {
				//断线开启重新链接
				d.conn.Close()
				log.Errorln(err.Error())
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
