package api

import (
	"go-dog/cache"
	"go-dog/cmd/define"
	"go-dog/internal/service"
	"go-dog/mysql"
	"go-dog/plugins"
	"go-dog/serviceinfo"
	"sync"

	"github.com/prometheus/common/log"
)

//APIService API服务
type _APIService struct {
	method  *serviceinfo.API
	name    string
	explain string
	count   int32
}

//Service 控制服务
type Service struct {
	service  plugins.Service
	mysql    *mysql.Mysql
	cache    *cache.Cache
	apis     map[string]*_APIService
	services map[string]*serviceinfo.APIServiceInfo
	lock     sync.RWMutex
}

//NewService 初始化服务
func NewService() *Service {
	ctl := new(Service)
	ctl.apis = make(map[string]*_APIService)
	ctl.services = make(map[string]*serviceinfo.APIServiceInfo)
	//初始化rpc服务端
	ctl.service = service.CreateService(define.TTL)
	//设置服务名称
	ctl.service.SetName(define.SvcController)
	//设置服务端最大访问量
	ctl.service.GetLimit().SetLimit(define.MaxServiceRequestCount)
	//设置客户端最大访问量
	ctl.service.GetClient().GetLimit().SetLimit(define.MaxClientRequestCount)
	//注册API上线通知
	ctl.service.GetClient().GetDiscovery().RegAPIServiceOnlineNotice(ctl._ApiServiceOnline)
	//注册API下线通知
	ctl.service.GetClient().GetDiscovery().RegAPIServiceOfflineNotice(ctl._ApiServiceOffline)
	//开始监听API事件
	ctl.service.GetClient().GetDiscovery().WatchAPIService()
	//初始化数据库
	ctl.mysql = mysql.NewMysql(ctl.service.GetCfg())
	//初始化API
	ctl.InitAPI()
	return ctl
}

//Run 启动
func (pointer *Service) Run() error {
	return pointer.service.Run()
}

//InitAPI 初始化API
func (pointer *Service) InitAPI() {
	//获取API列表
	pointer.service.RegisterAPI("GetAPIList", "v1", "get/api/list",
		plugins.POST,
		3,
		false,
		"获取api列表",
		pointer.GetAPIList)
	//获取服务列表
	pointer.service.RegisterAPI("GetServiceList", "v1", "get/service/list",
		plugins.POST,
		3,
		false,
		"获取服务列表",
		pointer.GetServiceList)
}

//apiServiceOnline api服务上线
func (pointer *Service) _ApiServiceOnline(key string, service *serviceinfo.APIServiceInfo) {
	pointer.lock.Lock()
	for _, method := range service.API {
		url := "/api/" + service.Name + "/" + method.Version + "/" + method.Path
		if api, ok := pointer.apis[url]; ok {
			api.count++
		} else {
			pointer.apis[url] = &_APIService{
				method:  method,
				name:    service.Name,
				explain: service.Explain,
				count:   1,
			}
		}
		pointer.services[key] = service
	}
	pointer.lock.Unlock()
}

//apiServiceOffline api服务下线
func (pointer *Service) _ApiServiceOffline(key string) {
	pointer.lock.Lock()
	if service, ok := pointer.services[key]; ok {
		for _, method := range service.API {
			url := "/api/" + service.Name + "/" + method.Version + "/" + method.Path
			if api, ok := pointer.apis[url]; ok {
				api.count--
				if api.count <= 0 {
					delete(pointer.apis, url)
				}
			}
		}
		delete(pointer.services, key)
	}
	pointer.lock.Unlock()
}

// Set 设置验证码ID
func (pointer *Service) Set(id string, value string) {
	if err := pointer.cache.GetCache().SetByTime(id, value, define.CodeValidityTime); err != nil {
		log.Error(err.Error())
	}
}

// Get 更具验证ID获取验证码
func (pointer *Service) Get(id string, clear bool) (vali string) {
	err := pointer.cache.GetCache().Get(id, &vali)
	if err != nil {
		log.Error(err.Error())
	}
	if clear {
		pointer.cache.GetCache().Del(id)
	}
	return
}

//Verify 验证验证码
func (pointer *Service) Verify(id, answer string, clear bool) bool {
	vali := ""
	err := pointer.cache.GetCache().Get(id, &vali)
	if err != nil {
		log.Error(err.Error())
	}
	if clear {
		pointer.cache.GetCache().Del(id)
	}
	if vali != answer {
		return false
	}
	return true
}
