package discovery

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/nacos"
	"github.com/tang-go/go-dog/plugins"
	"github.com/tang-go/go-dog/serviceinfo"
)

//Discovery 服务发现
type Discovery struct {
	ctx        context.Context
	cancel     context.CancelFunc
	cfg        plugins.Cfg
	watchAPI   bool
	watchRPC   bool
	closeheart chan bool
	apidata    map[string]*serviceinfo.APIServiceInfo
	rpcdata    map[string]*serviceinfo.RPCServiceInfo
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
		watchAPI:   false,
		watchRPC:   false,
		closeheart: make(chan bool),
		apidata:    make(map[string]*serviceinfo.APIServiceInfo),
		rpcdata:    make(map[string]*serviceinfo.RPCServiceInfo),
		apis:       make(map[string]*serviceinfo.ServcieAPI),
		gate:       "",
	}
	return dis
}

//WatchAPI 监听api服务--区分网关使用
func (d *Discovery) WatchAPI(gate string) {
	d.watchAPI = true
	d.gate = gate
	nacos.GetDiscovery().Discovery(
		d.ctx,
		"API",
		[]string{d.cfg.GetClusterName()},
		func(i nacos.Instance) {
			d.lock.Lock()
			defer d.lock.Unlock()
			info := new(serviceinfo.APIServiceInfo)
			info.Key = i.Metadata["Key"]
			info.Time = i.Metadata["Time"]
			info.Explain = i.Metadata["Explain"]
			info.Name = i.ServiceName
			info.Address = i.Ip
			info.Port = int(i.Port)
			if err := json.Unmarshal([]byte(i.Metadata["API"]), &info.API); err != nil {
				log.Traceln(err.Error())
				return
			}
			if err := json.Unmarshal([]byte(i.Metadata["Methods"]), &info.Methods); err != nil {
				log.Traceln(err.Error())
				return
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
					log.Tracef("api 上线 | %s | %s | %s ", info.Name, info.Key, url)
				}
			}
			info.API = apis
			d.apidata[info.Key] = info
			d.rpcdata[info.Key] = &serviceinfo.RPCServiceInfo{
				Key:       info.Key,
				Name:      info.Name,
				Address:   info.Address,
				Port:      info.Port,
				Methods:   info.Methods,
				Explain:   info.Explain,
				Longitude: info.Longitude,
				Latitude:  info.Latitude,
				Time:      info.Time,
			}
		}, func(i nacos.Instance) {
			d.lock.Lock()
			defer d.lock.Unlock()
			info := new(serviceinfo.APIServiceInfo)
			info.Key = i.Metadata["Key"]
			info.Time = i.Metadata["Time"]
			info.Explain = i.Metadata["Explain"]
			info.Name = i.ServiceName
			info.Address = i.Ip
			info.Port = int(i.Port)
			if err := json.Unmarshal([]byte(i.Metadata["API"]), &info.API); err != nil {
				log.Traceln(err.Error())
				return
			}
			apis := make([]*serviceinfo.API, 0)
			for _, method := range info.API {
				if method.Gate != d.gate {
					continue
				}
				apis = append(apis, method)
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
			delete(d.rpcdata, info.Key)
		},
	)
}

//WatchRPC 监听api服务
func (d *Discovery) WatchRPC() {
	d.watchRPC = true
	nacos.GetDiscovery().Discovery(
		d.ctx,
		"RPC",
		[]string{d.cfg.GetClusterName()},
		func(i nacos.Instance) {
			d.lock.Lock()
			defer d.lock.Unlock()
			info := new(serviceinfo.RPCServiceInfo)
			info.Key = i.Metadata["Key"]
			info.Time = i.Metadata["Time"]
			info.Explain = i.Metadata["Explain"]
			info.Name = i.ServiceName
			info.Address = i.Ip
			info.Port = int(i.Port)
			if err := json.Unmarshal([]byte(i.Metadata["Methods"]), &info.Methods); err != nil {
				log.Traceln(err.Error())
				return
			}
			d.rpcdata[info.Key] = info
			log.Tracef("rpc 上线 | %s | %s | %s:%d ", info.Name, info.Key, info.Address, info.Port)
		}, func(i nacos.Instance) {
			d.lock.Lock()
			defer d.lock.Unlock()
			info := new(serviceinfo.RPCServiceInfo)
			info.Key = i.Metadata["Key"]
			info.Time = i.Metadata["Time"]
			info.Explain = i.Metadata["Explain"]
			info.Name = i.ServiceName
			info.Address = i.Ip
			info.Port = int(i.Port)
			delete(d.rpcdata, info.Key)
			log.Tracef("rpc 下线 | %s | %s | %s:%d ", info.Name, info.Key, info.Address, info.Port)
		},
	)
}

//GetAllAPIService 获取所有API服务
func (d *Discovery) GetAllAPIService() (services []*serviceinfo.APIServiceInfo) {
	d.lock.RLock()
	for _, service := range d.apidata {
		services = append(services, service)
	}
	d.lock.RUnlock()
	return
}

//GetAllRPCService 获取所有RPC服务
func (d *Discovery) GetAllRPCService() (services []*serviceinfo.RPCServiceInfo) {
	d.lock.RLock()
	for _, service := range d.rpcdata {
		services = append(services, service)
	}
	d.lock.RUnlock()
	return
}

//GetRPCServiceByName 通过名称获取RPC服务
func (d *Discovery) GetRPCServiceByName(name string) (services []*serviceinfo.RPCServiceInfo) {
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
func (d *Discovery) GetAPIServiceByName(name string) (services []*serviceinfo.APIServiceInfo) {
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
