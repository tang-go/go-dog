package plugins

import (
	"context"

	"github.com/tang-go/go-dog/serviceinfo"
)

//Register 服务注册
type Register interface {

	//RegisterRPCService 注册RPC服务
	RegisterRPCService(ctx context.Context, info *serviceinfo.RPCServiceInfo) error

	//RegisterAPIService 注册API服务
	RegisterAPIService(ctx context.Context, info *serviceinfo.APIServiceInfo) error

	// Cancellation 注销服务
	Cancellation() error
}
