package plugins

import (
	"context"

	"github.com/tang-go/go-dog/serviceinfo"
)

//Register 服务注册
type Register interface {

	//RegisterRPCService 注册RPC服务
	RegisterRPCService(ctx context.Context, info *serviceinfo.ServiceInfo) error

	//RegisterHTTPService 注册HTTP服务
	RegisterHTTPService(ctx context.Context, info *serviceinfo.ServiceInfo) error

	// Cancellation 注销服务
	Cancellation() error
}
