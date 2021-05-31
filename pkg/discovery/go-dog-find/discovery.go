package discovery

import (
	"encoding/json"
	"net"
	"sync"
	"time"

	"github.com/tang-go/go-dog/lib/io"
	"github.com/tang-go/go-dog/lib/rand"
	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/pkg/discovery/go-dog-find/param"
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
	watchAPI   bool
	watchRPC   bool
	closeheart chan bool
	apidata    map[string]*serviceinfo.APIServiceInfo
	rpcdata    map[string]*serviceinfo.RPCServiceInfo
	apis       map[string]*serviceinfo.ServcieAPI
	gate       string
	lock       sync.RWMutex
}

//NewGoDogDiscovery  新建发现服务
func NewGoDogDiscovery(address []string) *GoDogDiscovery {
	dis := &GoDogDiscovery{
		address:    address,
		ttl:        2 * time.Second,
		count:      len(address),
		watchAPI:   false,
		watchRPC:   false,
		pos:        0,
		close:      false,
		closeheart: make(chan bool),
		apidata:    make(map[string]*serviceinfo.APIServiceInfo),
		rpcdata:    make(map[string]*serviceinfo.RPCServiceInfo),
		apis:       make(map[string]*serviceinfo.ServcieAPI),
		gate:       "",
	}
	return dis
}

//WatchAPI 监听api服务--区分网关使用
func (d *GoDogDiscovery) WatchAPI(gate string) {
	d.watchAPI = true
	d.gate = gate
}

//WatchRPC 监听api服务
func (d *GoDogDiscovery) WatchRPC() {
	d.watchRPC = true
}

//ConnectClient 建立链接
func (d *GoDogDiscovery) ConnectClient() {
	index := rand.IntRand(0, d.count)
	addr := d.address[index]
	if err := d._ConnectClient(addr); err != nil {
		panic(err)
	}
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

//GetAPIByURL 通过RUL获取API服务
func (d *GoDogDiscovery) GetAPIByURL(url string) (*serviceinfo.ServcieAPI, bool) {
	d.lock.RLock()
	defer d.lock.RUnlock()
	s, ok := d.apis[url]
	return s, ok
}

//RangeAPI 遍历api
func (d *GoDogDiscovery) RangeAPI(f func(url string, api *serviceinfo.ServcieAPI)) {
	d.lock.RLock()
	for url, api := range d.apis {
		f(url, api)
	}
	d.lock.RUnlock()
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

//_ConnectClient 建立链接
func (d *GoDogDiscovery) _ConnectClient(address string) error {
	d.lock.Lock()
	defer d.lock.Unlock()
	if d.close {
		return nil
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
		log.Traceln(err.Error())
		return err
	}
	//发送消息
	if err := d._SendMsg(conn, param.Login, buff); err != nil {
		//断线开启重新链接
		conn.Close()
		log.Traceln(err.Error())
		return err
	}
	d.conn = conn
	//开启心跳
	go d._Heart()
	//开启监听
	go d._Watch()
	//默认监听rpc服务消息
	if d.watchRPC {
		d._WatchRPCService()
	}
	//默认监听api服务消息
	if d.watchAPI {
		d._WatchAPIService()
	}
	log.Traceln("链接成功注册中心", address)
	return nil
}

//_SendMsg 发送消息
func (d *GoDogDiscovery) _SendMsg(conn net.Conn, cmd int8, buff []byte) error {
	event := new(param.Event)
	event.Cmd = cmd
	event.Data = buff
	data, err := event.EnCode(event)
	if err != nil {
		log.Traceln(err.Error())
		return err
	}
	if _, err := io.WriteByTime(conn, data, time.Now().Add(d.ttl)); err != nil {
		log.Traceln(err.Error())
		return err
	}
	return nil
}

//_Watch 开始监听
func (d *GoDogDiscovery) _Watch() {
	for {
		_, buff, err := io.Read(d.conn)
		if err != nil {
			d.closeheart <- true
			d.conn.Close()
			log.Traceln(err.Error())
			break
		}
		event := new(param.Event)
		if err := event.DeCode(buff, event); err != nil {
			log.Traceln(err.Error())
			continue
		}
		switch event.Cmd {
		//监听的上线消息
		case param.Listen:
			listen := new(param.ListenRes)
			if err := listen.DeCode(event.Data, listen); err != nil {
				log.Traceln(err.Error())
				continue
			}
			if listen.Label == param.APILabel {
				d._APIWatch(listen.Data)
			}
			if listen.Label == param.RPCLabel {
				d._RPCWatch(listen.Data)
			}
		}
	}
	for {
		for _, addr := range d.address {
			time.Sleep(d.ttl)
			log.Traceln("断线重链注册中心....")
			if d._ConnectClient(addr) == nil {
				return
			}
		}
	}
}

//_RPCWatch 推送事件
func (d *GoDogDiscovery) _RPCWatch(datas []param.Data) {
	d.lock.Lock()
	defer d.lock.Unlock()
	mp := make(map[string]string)
	for _, data := range datas {
		if _, ok := d.rpcdata[data.Key]; !ok {
			info := new(serviceinfo.RPCServiceInfo)
			if err := json.Unmarshal([]byte(data.Value), info); err != nil {
				log.Traceln(err.Error(), data.Key, data.Value)
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

//_APIWatch api推送事件
func (d *GoDogDiscovery) _APIWatch(datas []param.Data) {
	d.lock.Lock()
	defer d.lock.Unlock()
	mp := make(map[string]string)
	for _, data := range datas {
		if _, ok := d.apidata[data.Key]; !ok {
			info := new(serviceinfo.APIServiceInfo)
			if err := json.Unmarshal([]byte(data.Value), info); err != nil {
				log.Traceln(err.Error(), data.Key, data.Value)
				continue
			}
			apis := make([]*serviceinfo.API, 0)
			for _, method := range info.API {
				if method.Gate != d.gate {
					continue
				}
				apis = append(apis, method)
				url := method.Kind + method.Path
				if api, ok := d.apis[url]; ok {
					api.Count++
					d.apis[url] = api
				} else {
					d.apis[url] = &serviceinfo.ServcieAPI{
						Method:  method,
						Gate:    method.Gate,
						Tags:    method.Group,
						Explain: info.Explain,
						Name:    info.Name,
						Count:   1,
					}
					log.Tracef(" 上线 | %s | %s | %s ", info.Name, data.Key, url)
				}
			}
			info.API = apis
			d.apidata[data.Key] = info
		}
		mp[data.Key] = data.Value
	}
	for key, info := range d.apidata {
		if _, ok := mp[key]; !ok {
			for _, method := range info.API {
				url := method.Kind + method.Path
				if api, ok := d.apis[url]; ok {
					api.Count--
					if api.Count <= 0 {
						delete(d.apis, url)
						log.Tracef(" 下线 | %s | %s | %s ", info.Name, info.Key, url)
					}
				}
			}
			delete(d.apidata, key)
		}
	}
}

//_Heart 心跳
func (d *GoDogDiscovery) _Heart() {
	for {
		select {
		case <-d.closeheart:
			return
		case <-time.After(d.ttl):
			if err := d._SendMsg(d.conn, param.Heart, nil); err != nil {
				//断线开启重新链接
				d.conn.Close()
				log.Traceln(err.Error())
			}
		}

	}
}

//WatchRPCService 开始RPC服务发现
func (d *GoDogDiscovery) _WatchRPCService() {
	listen := &param.ListenReq{
		Label: param.RPCLabel,
	}
	buff, err := listen.EnCode(listen)
	if err != nil {
		panic(err.Error())
	}
	if err := d._SendMsg(d.conn, param.Listen, buff); err != nil {
		panic(err.Error())
	}
	log.Traceln("watch rpc")
}

//WatchAPIService 开始API服务发现
func (d *GoDogDiscovery) _WatchAPIService() {
	listen := &param.ListenReq{
		Label: param.APILabel,
	}
	buff, err := listen.EnCode(listen)
	if err != nil {
		panic(err.Error())
	}
	if err := d._SendMsg(d.conn, param.Listen, buff); err != nil {
		panic(err.Error())
	}
	log.Traceln("watch api")
}
