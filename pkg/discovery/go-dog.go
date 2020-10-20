package discovery

import (
	"encoding/json"
	"go-dog/cmd/go-dog-find/param"
	"go-dog/lib/io"
	"go-dog/log"
	"go-dog/serviceinfo"
	"net"
	"sync"
	"time"
)

//GoDogDiscovery 服务发现
type GoDogDiscovery struct {
	address                []string
	conn                   net.Conn
	ttl                    time.Duration
	pos                    int
	count                  int
	close                  bool
	closeheart             chan bool
	label                  map[string]string
	apidata                map[string]string
	rpcdata                map[string]string
	rpcServcieOnlineNotice func(string, *serviceinfo.ServiceInfo)
	rpcServcieOffineNotice func(string)
	apiServcieOnlineNotice func(string, *serviceinfo.APIServiceInfo)
	apiServcieOffineNotice func(string)
	lock                   sync.Mutex
}

//NewGoDogDiscovery  新建发现服务
func NewGoDogDiscovery(address []string, ttl int64) *GoDogDiscovery {
	dis := &GoDogDiscovery{
		address:    address,
		ttl:        2 * time.Second,
		count:      len(address),
		pos:        0,
		close:      false,
		closeheart: make(chan bool),
		label:      make(map[string]string),
		apidata:    make(map[string]string),
		rpcdata:    make(map[string]string),
	}
	if err := dis._ConnectClient(); err != nil {
		panic(err)
	}
	return dis
}

//_ConnectClient 建立链接
func (d *GoDogDiscovery) _ConnectClient() error {
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
	go d._Watch()
	log.Traceln("链接成功注册中心", address)
	return nil
}

//_Watch 开始监听
func (d *GoDogDiscovery) _Watch() {
	go d._Heart()
	if _, ok := d.label["/rpc"]; ok {
		d.WatchRPCService()
	}
	if _, ok := d.label["/api"]; ok {
		d.WatchAPIService()
	}
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
		if all.Label == "/rpc" {
			mp := make(map[string]string)
			for _, data := range all.Datas {
				if _, ok := d.rpcdata[data.Key]; !ok {
					if d.rpcServcieOnlineNotice != nil {
						info := new(serviceinfo.ServiceInfo)
						if err := json.Unmarshal([]byte(data.Value), info); err != nil {
							log.Errorln(err.Error(), data.Key, data.Value)
							continue
						}
						d.rpcServcieOnlineNotice(data.Key, info)
						d.rpcdata[data.Key] = data.Value
					}
				}
				mp[data.Key] = data.Value
			}
			for key := range d.rpcdata {
				if _, ok := mp[key]; !ok {
					if d.rpcServcieOffineNotice != nil {
						d.rpcServcieOffineNotice(key)
						delete(d.rpcdata, key)
					}
				}
			}
		}
		if all.Label == "/api" {
			mp := make(map[string]string)
			for _, data := range all.Datas {
				if _, ok := d.apidata[data.Key]; !ok {
					if d.apiServcieOnlineNotice != nil {
						info := new(serviceinfo.APIServiceInfo)
						if err := json.Unmarshal([]byte(data.Value), info); err != nil {
							log.Errorln(err.Error(), data.Key, data.Value)
							continue
						}
						d.apiServcieOnlineNotice(data.Key, info)
						d.apidata[data.Key] = data.Value
					}
				}
				mp[data.Key] = data.Value
			}
			for key := range d.apidata {
				if _, ok := mp[key]; !ok {
					if d.apiServcieOffineNotice != nil {
						d.apiServcieOffineNotice(key)
						delete(d.apidata, key)
					}
				}
			}
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

//RegRPCServiceOnlineNotice 注册RPC服务上线通知
func (d *GoDogDiscovery) RegRPCServiceOnlineNotice(f func(string, *serviceinfo.ServiceInfo)) {
	d.rpcServcieOnlineNotice = f
}

//RegRPCServiceOfflineNotice 注册RPC服务下线通知
func (d *GoDogDiscovery) RegRPCServiceOfflineNotice(f func(string)) {
	d.rpcServcieOffineNotice = f
}

//WatchRPCService 开始RPC服务发现
func (d *GoDogDiscovery) WatchRPCService() {
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
	d.label["/rpc"] = "/rpc"
	log.Traceln("watch /rpc")
}

//RegAPIServiceOnlineNotice 注册API服务上线通知
func (d *GoDogDiscovery) RegAPIServiceOnlineNotice(f func(string, *serviceinfo.APIServiceInfo)) {
	d.apiServcieOnlineNotice = f
}

//RegAPIServiceOfflineNotice 注册API服务下线通知
func (d *GoDogDiscovery) RegAPIServiceOfflineNotice(f func(string)) {
	d.apiServcieOffineNotice = f
}

//WatchAPIService 开始API服务发现
func (d *GoDogDiscovery) WatchAPIService() {
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
	d.label["/api"] = "/api"
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
