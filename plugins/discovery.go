package plugins

import (
	"github.com/tang-go/go-dog/serviceinfo"
)

//Discovery 服务发现
type Discovery interface {

	//GetAllAPIService 获取所有API服务
	GetAllAPIService() (services []*serviceinfo.APIServiceInfo)

	//GetAllRPCService 获取所有RPC服务
	GetAllRPCService() (services []*serviceinfo.RPCServiceInfo)

	//GetRPCServiceByName 通过名称获取RPC服务
	GetRPCServiceByName(name string) (services []*serviceinfo.RPCServiceInfo)

	//GetAPIServiceByName 通过名称获取API服务
	GetAPIServiceByName(name string) (services []*serviceinfo.APIServiceInfo)

	//Close 关闭服务
	Close() error
}
