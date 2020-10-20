package plugins

import (
	"github.com/tang-go/go-dog/serviceinfo"
)

//Discovery 服务发现
type Discovery interface {

	//RegRPCServiceOnlineNotice 注册RPC服务上线通知
	RegRPCServiceOnlineNotice(f func(string, *serviceinfo.ServiceInfo))

	//RegAPIServiceOfflineNotice 注册RPC服务下线通知
	RegAPIServiceOfflineNotice(f func(string))

	//RegAPIServiceOnlineNotice 注册API服务上线通知
	RegAPIServiceOnlineNotice(f func(string, *serviceinfo.APIServiceInfo))

	//RegRPCServiceOfflineNotice 注册API服务下线通知
	RegRPCServiceOfflineNotice(f func(string))

	//WatchRPCService 开始RPC服务发现
	WatchRPCService()

	//WatchAPIService 开始API服务发现
	WatchAPIService()

	//Close 关闭服务
	Close() error
}
