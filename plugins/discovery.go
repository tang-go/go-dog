package plugins

import (
	"github.com/tang-go/go-dog/serviceinfo"
)

//Discovery 服务发现
type Discovery interface {

	//WatchAPI 监听api服务--区分网关使用
	WatchAPI(gate string)

	//WatchRPC 监听api服务
	WatchRPC()

	//GetAllAPIService 获取所有API服务
	GetAllAPIService() (services []*serviceinfo.ServiceInfo)

	//GetAllRPCService 获取所有RPC服务
	GetAllRPCService() (services []*serviceinfo.ServiceInfo)

	//GetRPCServiceByName 通过名称获取RPC服务
	GetRPCServiceByName(name string) (services []*serviceinfo.ServiceInfo)

	//GetAPIServiceByName 通过名称获取API服务
	GetAPIServiceByName(name string) (services []*serviceinfo.ServiceInfo)

	//GetAPIByURL 通过RUL获取API服务
	GetAPIByURL(url string) (*serviceinfo.ServcieAPI, bool)

	//RangeAPI 遍历api
	RangeAPI(f func(url string, api *serviceinfo.ServcieAPI))

	//Close 关闭服务
	Close() error
}
