package register

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/tang-go/go-dog/nacos"
	"github.com/tang-go/go-dog/plugins"
	"github.com/tang-go/go-dog/serviceinfo"
)

//EtcdRegister nacos 服务注册
type Register struct {
	cfg plugins.Cfg
}

//NewNacosRegister 初始化一个nacos服务注册中心
func NewNacosRegister(cfg plugins.Cfg) *Register {
	r := new(Register)
	r.cfg = cfg
	return r
}

//RegisterRPCService 注册RPC服务
func (s *Register) RegisterRPCService(ctx context.Context, info *serviceinfo.RPCServiceInfo) error {
	key := "rpc/" + fmt.Sprintf("%s:%d", info.Address, info.Port)
	info.Key = key
	methods, _ := json.Marshal(info.Methods)
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
			"Methods":   string(methods),
			"Time":      info.Time,
			"Key":       info.Key,
			"Name":      info.Name,
			"Longitude": fmt.Sprintf("%d", info.Longitude),
			"Latitude":  fmt.Sprintf("%d", info.Latitude),
			"Explain":   info.Explain,
		},
	}
	nacos.GetRegister().Register(param)
	return nil
}

//RegisterAPIService 注册API服务
func (s *Register) RegisterAPIService(ctx context.Context, info *serviceinfo.APIServiceInfo) error {
	key := "api/" + fmt.Sprintf("%s:%d", info.Address, info.Port)
	info.Key = key
	api, _ := json.Marshal(info.API)
	methods, _ := json.Marshal(info.Methods)
	param := nacos.RegisterInstanceParam{
		Ip:          info.Address,
		Port:        uint64(info.Port),
		Weight:      10,
		Enable:      true,
		Healthy:     true,
		ClusterName: s.cfg.GetClusterName(),
		ServiceName: info.Name,
		GroupName:   "API",
		Ephemeral:   true,
		Metadata: map[string]string{
			"Methods":   string(methods),
			"API":       string(api),
			"Time":      info.Time,
			"Key":       info.Key,
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
