package plugins

import (
	"context"
	"go-dog/serviceinfo"
)

//Register 服务注册
type Register interface {

	//RegisterRPCService 注册RPC服务
	RegisterRPCService(ctx context.Context, info *serviceinfo.ServiceInfo) error

	//RegisterAPIService 注册API服务
	RegisterAPIService(ctx context.Context, info *serviceinfo.APIServiceInfo) error

	// Cancellation 注销服务
	Cancellation() error
}
