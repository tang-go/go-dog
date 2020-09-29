package controller

import (
	"go-dog/cmd/go-dog-ctl/param"
	"go-dog/define"
	"go-dog/internal/client"
	"go-dog/internal/config"
	"go-dog/internal/discovery"
	"go-dog/internal/service"
	"go-dog/plugins"
	"go-dog/serviceinfo"
	"sync"
)

//APIService API服务
type APIService struct {
	method  *serviceinfo.API
	name    string
	explain string
	count   int32
}

//Controller 控制服务
type Controller struct {
	service   plugins.Service
	cfg       plugins.Cfg
	discovery plugins.Discovery
	client    plugins.Client
	apis      map[string]*APIService
	services  map[string]*serviceinfo.APIServiceInfo
	lock      sync.RWMutex
}

//NewController 初始化服务
func NewController() *Controller {
	ctl := new(Controller)
	ctl.apis = make(map[string]*APIService)
	ctl.services = make(map[string]*serviceinfo.APIServiceInfo)
	ctl.cfg = config.NewConfig()
	ctl.discovery = discovery.NewEtcdDiscovery(ctl.cfg.GetEtcd(), define.TTL)
	ctl.client = client.NewClient(ctl.cfg, ctl.discovery)
	ctl.discovery.RegAPIServiceOnlineNotice(ctl.apiServiceOnline)
	ctl.discovery.RegAPIServiceOfflineNotice(ctl.apiServiceOffline)
	ctl.discovery.WatchAPIService()
	ctl.service = service.CreateService(ctl.cfg, ctl.client)
	ctl.Init()
	return ctl
}

//Run 启动
func (c *Controller) Run() error {
	return c.service.Run()
}

//Init 初始化
func (c *Controller) Init() {
	c.service.RegisterAPI("GetAPIList", "v1", "get/api/list", plugins.POST, 3, false, "获取api列表", c.GetAPIList)
	c.service.RegisterAPI("GetServiceList", "v1", "get/service/list", plugins.POST, 3, false, "获取服务列表", c.GetServiceList)
}

//GetServiceList 获取服务列表
func (c *Controller) GetServiceList(ctx plugins.Context, req param.GetServiceReq) (res param.GetServiceRes, err error) {
	services := c.service.GetClient().GetAllService()
	for _, service := range services {
		s := &param.ServiceInfo{
			Key:       service.Key,
			Name:      service.Name,
			Address:   service.Address,
			Port:      service.Port,
			Explain:   service.Explain,
			Longitude: service.Longitude,
			Latitude:  service.Latitude,
			Time:      service.Time,
		}
		for _, method := range service.Methods {
			s.Methods = append(s.Methods, &param.Method{
				Name:     method.Name,
				Level:    method.Level,
				Request:  method.Request,
				Response: method.Response,
				Explain:  method.Explain,
				IsAuth:   method.IsAuth,
			})
		}
		res.List = append(res.List, s)
	}
	return
}

//GetAPIList 获取api列表
func (c *Controller) GetAPIList(ctx plugins.Context, req param.GetAPIListReq) (res param.GetAPIListRes, err error) {
	list := make(map[string]*param.Service)
	c.lock.RLock()
	for key, api := range c.apis {
		if service, ok := list[api.name]; ok {
			service.APIS = append(service.APIS, &param.API{
				Name:     api.method.Name,
				Level:    api.method.Level,
				Request:  api.method.Request,
				Response: api.method.Response,
				Explain:  api.method.Explain,
				IsAuth:   api.method.IsAuth,
				Version:  api.method.Version,
				URL:      key,
				Kind:     api.method.Kind,
			})
		} else {
			s := &param.Service{
				Name:    api.name,
				Explain: api.explain,
				APIS: []*param.API{
					&param.API{
						Name:     api.method.Name,
						Level:    api.method.Level,
						Request:  api.method.Request,
						Response: api.method.Response,
						Explain:  api.method.Explain,
						IsAuth:   api.method.IsAuth,
						Version:  api.method.Version,
						URL:      key,
						Kind:     api.method.Kind,
					},
				},
			}
			list[api.name] = s
		}
	}
	c.lock.RUnlock()
	for _, s := range list {
		res.List = append(res.List, s)
	}
	return
}

//apiServiceOnline api服务上线
func (c *Controller) apiServiceOnline(key string, service *serviceinfo.APIServiceInfo) {
	c.lock.Lock()
	for _, method := range service.API {
		url := "/api/" + service.Name + "/" + method.Version + "/" + method.Path
		if api, ok := c.apis[url]; ok {
			api.count++
		} else {
			c.apis[url] = &APIService{
				method:  method,
				name:    service.Name,
				explain: service.Explain,
				count:   1,
			}
		}
		c.services[key] = service
	}
	c.lock.Unlock()
}

//apiServiceOffline api服务下线
func (c *Controller) apiServiceOffline(key string) {
	c.lock.Lock()
	if service, ok := c.services[key]; ok {
		for _, method := range service.API {
			url := "/api/" + service.Name + "/" + method.Version + "/" + method.Path
			if api, ok := c.apis[url]; ok {
				api.count--
				if api.count <= 0 {
					delete(c.apis, url)
				}
			}
		}
		delete(c.services, key)
	}
	c.lock.Unlock()
}
