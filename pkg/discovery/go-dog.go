package discovery

import (
	"encoding/json"
	"net"
	"sync"
	"time"

	"github.com/tang-go/go-dog/lib/io"
	"github.com/tang-go/go-dog/lib/rand"
	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/pkg/discovery/param"
	"github.com/tang-go/go-dog/serviceinfo"
)

//GoDogDiscovery 服务发现
type GoDogDiscovery struct {
	address    []string
	conn       net.Conn
	ttl        time.Duration
	pos        int
	count      int
	close      bool
	closeheart chan bool
	apidata    map[string]*serviceinfo.APIServiceInfo
	rpcdata    map[string]*serviceinfo.RPCServiceInfo
	lock       sync.RWMutex
}

//NewGoDogDiscovery  新建发现服务
func NewGoDogDiscovery(address []string) *GoDogDiscovery {
	dis := &GoDogDiscovery{
		address:    address,
		ttl:        2 * time.Second,
		count:      len(address),
		pos:        0,
		close:      false,
		closeheart: make(chan bool),
		apidata:    make(map[string]*serviceinfo.APIServiceInfo),
		rpcdata:    make(map[string]*serviceinfo.RPCServiceInfo),
	}
	//初始化第一个链接
	if err := dis._ConnectClient(); err != nil {
		panic(err)
	}
	//等待一个心跳时间
	time.Sleep(dis.ttl)
	return dis
}

//GetAllAPIService 获取所有API服务
func (d *GoDogDiscovery) GetAllAPIService() (services []*serviceinfo.APIServiceInfo) {
	d.lock.RLock()
	for _, service := range d.apidata {
		services = append(services, service)
	}
	d.lock.RUnlock()
	return
}

//GetAllRPCService 获取所有RPC服务
func (d *GoDogDiscovery) GetAllRPCService() (services []*serviceinfo.RPCServiceInfo) {
	d.lock.RLock()
	for _, service := range d.rpcdata {
		services = append(services, service)
	}
	d.lock.RUnlock()
	return
}

//GetRPCServiceByName 通过名称获取RPC服务
func (d *GoDogDiscovery) GetRPCServiceByName(name string) (services []*serviceinfo.RPCServiceInfo) {
	d.lock.RLock()
	for _, service := range d.rpcdata {
		if service.Name == name {
			services = append(services, service)
		}
	}
	d.lock.RUnlock()
	return
}

//GetAPIServiceByName 通过名称获取API服务
func (d *GoDogDiscovery) GetAPIServiceByName(name string) (services []*serviceinfo.APIServiceInfo) {
	d.lock.RLock()
	for _, service := range d.apidata {
		if service.Name == name {
			services = append(services, service)
		}
	}
	d.lock.RUnlock()
	return
}

//_ConnectClient 建立链接
func (d *GoDogDiscovery) _ConnectClient() error {
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
	login.Type = param.DisType
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
	//开启心跳
	go d._Heart()
	//开启监听
	go d._Watch()
	//默认监听rpc服务消息
	d._WatchRPCService()
	log.Traceln("链接成功注册中心", address)
	return nil
}

//_Watch 开始监听
func (d *GoDogDiscovery) _Watch() {
	for {
		_, buff, err := io.Read(d.conn)
		if err != nil {
			d.closeheart <- true
			d.conn.Close()
			log.Errorln(err.Error())
			break
		}
		all := new(param.All)
		if err := all.DeCode(buff, all); err != nil {
			log.Errorln(err.Error())
			continue
		}
		d.lock.Lock()
		if all.Label == "/rpc" {
			mp := make(map[string]string)
			for _, data := range all.Datas {
				if _, ok := d.rpcdata[data.Key]; !ok {
					info := new(serviceinfo.RPCServiceInfo)
					if err := json.Unmarshal([]byte(data.Value), info); err != nil {
						log.Errorln(err.Error(), data.Key, data.Value)
						continue
					}
					d.rpcdata[data.Key] = info
					log.Tracef("rpc 上线 | %s | %s | %s ", info.Name, data.Key, info.Address)
				}
				mp[data.Key] = data.Value
			}
			for key, info := range d.rpcdata {
				if _, ok := mp[key]; !ok {
					delete(d.rpcdata, key)
					log.Tracef("rpc 下线 | %s | %s | %s ", info.Name, key, info.Address)
				}
			}
		}
		if all.Label == "/api" {
			mp := make(map[string]string)
			for _, data := range all.Datas {
				if _, ok := d.apidata[data.Key]; !ok {
					info := new(serviceinfo.APIServiceInfo)
					if err := json.Unmarshal([]byte(data.Value), info); err != nil {
						log.Errorln(err.Error(), data.Key, data.Value)
						continue
					}
					d.apidata[data.Key] = info
					log.Tracef("api 上线 | %s | %s | %s ", info.Name, data.Key, info.Address)
				}
				mp[data.Key] = data.Value
			}
			for key, info := range d.apidata {
				if _, ok := mp[key]; !ok {
					delete(d.apidata, key)
					log.Tracef("api 下线 | %s | %s | %s ", info.Name, key, info.Address)
				}
			}
		}
		d.lock.Unlock()
	}

	for {
		time.Sleep(d.ttl)
		log.Traceln("断线重链注册中心....")
		if d._ConnectClient() == nil {
			return
		}
	}
}

//_Heart 心跳
func (d *GoDogDiscovery) _Heart() {
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

//WatchRPCService 开始RPC服务发现
func (d *GoDogDiscovery) _WatchRPCService() {
	//开启监听
	listen := &param.Event{
		Cmd:   param.Listen,
		Label: "/rpc",
	}
	buff, err := listen.EnCode(listen)
	if err != nil {
		panic(err.Error())
	}
	if _, err := io.WriteByTime(d.conn, buff, time.Now().Add(d.ttl)); err != nil {
		panic(err.Error())
	}
	log.Traceln("watch /rpc")
}

//WatchAPIService 开始API服务发现
func (d *GoDogDiscovery) _WatchAPIService() {
	//开启监听
	listen := &param.Event{
		Cmd:   param.Listen,
		Label: "/api",
	}
	buff, err := listen.EnCode(listen)
	if err != nil {
		panic(err.Error())
	}
	if _, err := io.WriteByTime(d.conn, buff, time.Now().Add(d.ttl)); err != nil {
		panic(err.Error())
	}
	log.Traceln("watch /api")
}

//Close 关闭服务
func (d *GoDogDiscovery) Close() error {
	d.lock.Lock()
	defer d.lock.Unlock()
	d.close = true
	if d.conn != nil {
		return d.conn.Close()
	}
	return nil
}
