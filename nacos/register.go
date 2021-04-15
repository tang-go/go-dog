package nacos

import (
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

//Register 注册中心
type Register struct {
	client  naming_client.INamingClient
	servers []vo.RegisterInstanceParam
}

//NewRegister 创建注册中心
func newRegister(c naming_client.INamingClient) *Register {
	r := new(Register)
	r.client = c
	return r
}

//Register 注册一个服务
func (r *Register) Register(info vo.RegisterInstanceParam) (bool, error) {
	r.servers = append(r.servers, info)
	return r.client.RegisterInstance(info)
}

//DeregisterInstance 取消注册一个服务
func (r *Register) DeregisterInstance() {
	for _, info := range r.servers {
		r.client.DeregisterInstance(vo.DeregisterInstanceParam{
			Ip:          info.Ip,
			Port:        info.Port,
			ServiceName: info.ServiceName,
			Ephemeral:   info.Ephemeral,
			Cluster:     info.ClusterName,
			GroupName:   info.GroupName,
		})
	}
}
