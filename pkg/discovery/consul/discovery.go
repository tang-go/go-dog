package discovery

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/tang-go/go-dog/consul"
	"github.com/tang-go/go-dog/lib/net"
	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/plugins"
	"github.com/tang-go/go-dog/serviceinfo"
)

var apiOnce sync.Once
var rpcOnce sync.Once

//Discovery 服务发现
type Discovery struct {
	ctx        context.Context
	cancel     context.CancelFunc
	cfg        plugins.Cfg
	closeheart chan bool
	apidata    map[string]*serviceinfo.ServiceInfo
	rpcdata    map[string]*serviceinfo.ServiceInfo
	apis       map[string]*serviceinfo.ServcieAPI
	gate       string
	lock       sync.RWMutex
}

//NewDiscovery  新建发现服务
func NewDiscovery(cfg plugins.Cfg) *Discovery {
	ctx, cancel := context.WithCancel(context.Background())
	dis := &Discovery{
		ctx:        ctx,
		cancel:     cancel,
		cfg:        cfg,
		closeheart: make(chan bool),
		apidata:    make(map[string]*serviceinfo.ServiceInfo),
		rpcdata:    make(map[string]*serviceinfo.ServiceInfo),
		apis:       make(map[string]*serviceinfo.ServcieAPI),
		gate:       "",
	}
	consul.Init(cfg.GetConsul())
	dis.WatchRPC()
	return dis
}

//WatchAPI 监听api服务--区分网关使用
func (d *Discovery) WatchAPI(gate string) {
	apiOnce.Do(func() {
		log.Traceln("监听api")
		d.watchAPI(gate)
	})
}

//WatchRPC 监听rpc服务
func (d *Discovery) WatchRPC() {
	rpcOnce.Do(func() {
		log.Traceln("监听rpc")
		d.watchRPC()
	})
}

//WatchAPI 监听api服务--区分网关使用
func (d *Discovery) watchAPI(gate string) {
	d.gate = gate
	consul.GetDiscovery().Discovery(
		d.ctx,
		[]string{"HTTP", d.cfg.GetClusterName()},
		func(i consul.Instance) {
			d.lock.Lock()
			defer d.lock.Unlock()
			key := fmt.Sprintf("%s:%d", i.Address, i.Port)
			if _, ok := d.apidata[key]; ok {
				log.Traceln(key, "已经存在")
				return
			}
			url := fmt.Sprintf("http://%s:%d/apis", i.Address, i.Port)
			apiConfig, err := net.HttpsGet(url)
			if err != nil {
				log.Errorln(err.Error())
				return
			}
			info := new(serviceinfo.ServiceInfo)
			if err := json.Unmarshal(apiConfig, info); err != nil {
				log.Errorln(err.Error())
				return
			}
			info.Key = key
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
					log.Tracef("api 上线 | %s | %s | %s ", info.Name, info.Key, url)
				}
			}
			d.apidata[info.Key] = info
		}, func(i consul.Instance) {
			d.lock.Lock()
			defer d.lock.Unlock()
			key := fmt.Sprintf("%s:%d", i.Address, i.Port)
			info, ok := d.apidata[key]
			if !ok {
				log.Traceln(key, "不存在")
				return
			}
			for _, method := range info.API {
				if method.Gate != d.gate {
					continue
				}
				url := method.Kind + method.Path
				if api, ok := d.apis[url]; ok {
					api.Count--
					if api.Count <= 0 {
						delete(d.apis, url)
						log.Tracef("api 下线 | %s | %s | %s ", info.Name, info.Key, url)
					}
				}
			}
			delete(d.apidata, info.Key)
		},
	)
}

//WatchRPC 监听api服务
func (d *Discovery) watchRPC() {
	consul.GetDiscovery().Discovery(
		d.ctx,
		[]string{"RPC", d.cfg.GetClusterName()},
		func(i consul.Instance) {
			d.lock.Lock()
			defer d.lock.Unlock()
			info := new(serviceinfo.ServiceInfo)
			info.Group = "RPC"
			info.Time = i.Meta["Time"]
			info.Explain = i.Meta["Explain"]
			info.Name = i.Service
			info.Address = i.Address
			info.Port = int(i.Port)
			info.Key = fmt.Sprintf("%s:%d", info.Address, info.Port)
			d.rpcdata[info.Key] = info
			log.Tracef("rpc 上线 | %s | %s | %s:%d ", info.Name, info.Key, info.Address, info.Port)
		}, func(i consul.Instance) {
			d.lock.Lock()
			defer d.lock.Unlock()
			info := new(serviceinfo.ServiceInfo)
			info.Group = "RPC"
			info.Time = i.Meta["Time"]
			info.Explain = i.Meta["Explain"]
			info.Name = i.Service
			info.Address = i.Address
			info.Port = int(i.Port)
			info.Key = fmt.Sprintf("%s:%d", info.Address, info.Port)
			delete(d.rpcdata, info.Key)
			log.Tracef("rpc 下线 | %s | %s | %s:%d ", info.Name, info.Key, info.Address, info.Port)
		},
	)
}

//GetRPCServiceByName 通过名称获取RPC服务
func (d *Discovery) GetRPCServiceByName(name string) (services []*serviceinfo.ServiceInfo) {
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
func (d *Discovery) GetAPIByURL(url string) (*serviceinfo.ServcieAPI, bool) {
	d.lock.RLock()
	defer d.lock.RUnlock()
	s, ok := d.apis[url]
	return s, ok
}

//RangeAPI 遍历api
func (d *Discovery) RangeAPI(f func(url string, api *serviceinfo.ServcieAPI)) {
	d.lock.RLock()
	for url, api := range d.apis {
		f(url, api)
	}
	d.lock.RUnlock()
}

//GetAPIServiceByName 通过名称获取API服务
func (d *Discovery) GetAPIServiceByName(name string) (services []*serviceinfo.ServiceInfo) {
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
func (d *Discovery) Close() error {
	d.cancel()
	return nil
}
