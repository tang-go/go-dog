package register

import (
	"context"
	"fmt"

	"github.com/tang-go/go-dog/consul"
	"github.com/tang-go/go-dog/plugins"
	"github.com/tang-go/go-dog/serviceinfo"
)

//EtcdRegister consul 服务注册
type Register struct {
	cfg    plugins.Cfg
	dataID string
	group  string
}

//NewConsulRegister 初始化一个consul服务注册中心
func NewConsulRegister(cfg plugins.Cfg) *Register {
	r := new(Register)
	r.cfg = cfg
	consul.Init(cfg.GetConsul())
	return r
}

//RegisterRPCService 注册RPC服务
func (s *Register) RegisterRPCService(ctx context.Context, info *serviceinfo.ServiceInfo) error {
	info.Key = fmt.Sprintf("%s:%d", info.Address, info.Port)
	info.Group = "RPC"
	param := consul.RegisterInstanceParam{
		Address: info.Address,
		Port:    info.Port,
		Tags:    []string{s.cfg.GetClusterName(), "RPC"},
		Name:    info.Name,
		Meta: map[string]string{
			"Time":      info.Time,
			"Name":      info.Name,
			"Longitude": fmt.Sprintf("%d", info.Longitude),
			"Latitude":  fmt.Sprintf("%d", info.Latitude),
			"Explain":   info.Explain,
		},
	}
	consul.GetRegister().Register(param)
	return nil
}

//RegisterHTTPService 注册HTTP服务
func (s *Register) RegisterHTTPService(ctx context.Context, info *serviceinfo.ServiceInfo) error {
	info.Key = fmt.Sprintf("%s:%d", info.Address, info.Port)
	info.Group = "HTTP"
	param := consul.RegisterInstanceParam{
		Address: info.Address,
		Port:    info.Port,
		Tags:    []string{s.cfg.GetClusterName(), "HTTP"},
		Name:    info.Name,
		Meta: map[string]string{
			"Time":      info.Time,
			"Name":      info.Name,
			"Longitude": fmt.Sprintf("%d", info.Longitude),
			"Latitude":  fmt.Sprintf("%d", info.Latitude),
			"Explain":   info.Explain,
		},
	}
	consul.GetRegister().Register(param)
	return nil
}

// Cancellation 注销服务
func (s *Register) Cancellation() error {
	consul.GetRegister().DeregisterInstance()
	return nil
}
