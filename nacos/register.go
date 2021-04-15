package nacos

import (
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

//注册参数
type RegisterInstanceParam struct {
	Ip          string            `param:"ip"`          //required
	Port        uint64            `param:"port"`        //required
	Weight      float64           `param:"weight"`      //required,it must be lager than 0
	Enable      bool              `param:"enabled"`     //required,the instance can be access or not
	Healthy     bool              `param:"healthy"`     //required,the instance is health or not
	Metadata    map[string]string `param:"metadata"`    //optional
	ClusterName string            `param:"clusterName"` //optional,default:DEFAULT
	ServiceName string            `param:"serviceName"` //required
	GroupName   string            `param:"groupName"`   //optional,default:DEFAULT_GROUP
	Ephemeral   bool              `param:"ephemeral"`   //optional
}

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
func (r *Register) Register(info RegisterInstanceParam) (bool, error) {
	param := vo.RegisterInstanceParam{
		Ip:          info.Ip,
		Port:        info.Port,
		Weight:      info.Weight,
		Enable:      info.Enable,
		Healthy:     info.Healthy,
		Metadata:    info.Metadata,
		ClusterName: info.ClusterName,
		ServiceName: info.ServiceName,
		GroupName:   info.GroupName,
		Ephemeral:   info.Ephemeral,
	}
	r.servers = append(r.servers, param)
	return r.client.RegisterInstance(param)
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
