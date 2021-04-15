package nacos

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

//Instance 实例内容
type Instance struct {
	Valid       bool              `json:"valid"`
	Marked      bool              `json:"marked"`
	InstanceId  string            `json:"instanceId"`
	Port        uint64            `json:"port"`
	Ip          string            `json:"ip"`
	Weight      float64           `json:"weight"`
	Metadata    map[string]string `json:"metadata"`
	ClusterName string            `json:"clusterName"`
	ServiceName string            `json:"serviceName"`
	Enable      bool              `json:"enabled"`
	Healthy     bool              `json:"healthy"`
	Ephemeral   bool              `json:"ephemeral"`
}

//Discovery 服务发现
type Discovery struct {
	client naming_client.INamingClient
}

//NewDiscovery 创建一个服务发现
func newDiscovery(c naming_client.INamingClient) *Discovery {
	d := new(Discovery)
	d.client = c
	return d
}

//Discovery 服务发现
func (d *Discovery) Discovery(ctx context.Context, groupName string, clusters []string, up func(Instance), down func(Instance)) {
	go func() {
		var listen sync.Map
		for {
			select {
			case <-ctx.Done():
				return
			case <-time.After(time.Second):
				var page uint32 = 1
				var count int64 = 0
				var size uint32 = 100
				serviceInfo := make(map[string]model.Instance)
				for {
					serviceInfos, err := d.client.GetAllServicesInfo(vo.GetAllServiceInfoParam{
						GroupName: groupName,
						PageNo:    page,
						PageSize:  size,
					})
					if err != nil {
						panic(err)
					}
					if len(serviceInfos.Doms) <= 0 || serviceInfos.Count <= 0 {
						break
					}

					for _, name := range serviceInfos.Doms {
						instances, err := d.client.SelectAllInstances(vo.SelectAllInstancesParam{
							Clusters:    clusters,
							GroupName:   groupName,
							ServiceName: name,
						})
						if err != nil {
							fmt.Println(err.Error())
							continue
						}
						for _, instance := range instances {
							if instance.Healthy != true {
								continue
							}
							instance.ServiceName = name
							serviceInfo[fmt.Sprintf("%s:%d", instance.Ip, instance.Port)] = instance
						}
					}
					count = count + 1
					page++
				}
				for key, value := range serviceInfo {
					if _, ok := listen.Load(key); !ok {
						up(Instance{
							Valid:       value.Valid,
							Marked:      value.Marked,
							InstanceId:  value.InstanceId,
							Port:        value.Port,
							Ip:          value.Ip,
							Weight:      value.Weight,
							Metadata:    value.Metadata,
							ClusterName: value.ClusterName,
							ServiceName: value.ServiceName,
							Enable:      value.Enable,
							Healthy:     value.Healthy,
							Ephemeral:   value.Ephemeral,
						})
					}
					listen.Store(key, value)
				}
				listen.Range(func(key, info interface{}) bool {
					if _, ok := serviceInfo[key.(string)]; !ok {
						listen.Delete(key)
						value := info.(model.Instance)
						down(Instance{
							Valid:       value.Valid,
							Marked:      value.Marked,
							InstanceId:  value.InstanceId,
							Port:        value.Port,
							Ip:          value.Ip,
							Weight:      value.Weight,
							Metadata:    value.Metadata,
							ClusterName: value.ClusterName,
							ServiceName: value.ServiceName,
							Enable:      value.Enable,
							Healthy:     value.Healthy,
							Ephemeral:   value.Ephemeral,
						})
					}
					return true
				})
			}
		}
	}()
}
