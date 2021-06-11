package register

import (
	"context"
	"fmt"

	"github.com/tang-go/go-dog/nacos"
	"github.com/tang-go/go-dog/plugins"
	"github.com/tang-go/go-dog/serviceinfo"
)

//EtcdRegister nacos 服务注册
type Register struct {
	cfg    plugins.Cfg
	dataID string
	group  string
}

//NewNacosRegister 初始化一个nacos服务注册中心
func NewNacosRegister(cfg plugins.Cfg) *Register {
	r := new(Register)
	r.cfg = cfg
	return r
}

//RegisterRPCService 注册RPC服务
func (s *Register) RegisterRPCService(ctx context.Context, info *serviceinfo.ServiceInfo) error {
	info.Key = fmt.Sprintf("%s:%d", info.Address, info.Port)
	info.Group = "RPC"
	param := nacos.RegisterInstanceParam{
		Ip:          info.Address,
		Port:        uint64(info.Port),
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		ClusterName: s.cfg.GetClusterName(),
		ServiceName: info.Name,
		GroupName:   "RPC",
		Ephemeral:   true,
		Metadata: map[string]string{
			"Time":      info.Time,
			"Name":      info.Name,
			"Longitude": fmt.Sprintf("%d", info.Longitude),
			"Latitude":  fmt.Sprintf("%d", info.Latitude),
			"Explain":   info.Explain,
		},
	}
	nacos.GetRegister().Register(param)
	return nil
}

//RegisterHTTPService 注册HTTP服务
func (s *Register) RegisterHTTPService(ctx context.Context, info *serviceinfo.ServiceInfo) error {
	info.Key = fmt.Sprintf("%s:%d", info.Address, info.Port)
	info.Group = "HTTP"
	param := nacos.RegisterInstanceParam{
		Ip:          info.Address,
		Port:        uint64(info.Port),
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		ClusterName: s.cfg.GetClusterName(),
		ServiceName: info.Name,
		GroupName:   "HTTP",
		Ephemeral:   true,
		Metadata: map[string]string{
			"Time":      info.Time,
			"Name":      info.Name,
			"Longitude": fmt.Sprintf("%d", info.Longitude),
			"Latitude":  fmt.Sprintf("%d", info.Latitude),
			"Explain":   info.Explain,
		},
	}
	nacos.GetRegister().Register(param)
	return nil
}

// Cancellation 注销服务
func (s *Register) Cancellation() error {
	nacos.GetRegister().DeregisterInstance()
	return nil
}
