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

//Discovery 服务发现
type Discovery struct {
	client naming_client.INamingClient
	listen map[string]*sync.Map
	lock   sync.RWMutex
}

//NewDiscovery 创建一个服务发现
func newDiscovery(c naming_client.INamingClient) *Discovery {
	d := new(Discovery)
	d.client = c
	d.listen = make(map[string]*sync.Map)
	return d
}

//Discovery 服务发现
func (d *Discovery) Discovery(ctx context.Context, groupName string, clusters []string, up func(model.Instance), down func(model.Instance)) {
	d.lock.Lock()
	d.listen[groupName] = new(sync.Map)
	d.lock.Unlock()
	go func() {
		for {
			select {
			case <-ctx.Done():
				d.lock.Lock()
				delete(d.listen, groupName)
				d.lock.Unlock()
				return
			case <-time.After(time.Second):
				d.discovery(groupName, clusters, up, down)
			}
		}
	}()
}

func (d *Discovery) discovery(groupName string, clusters []string, up func(model.Instance), down func(model.Instance)) {
	var page uint32 = 1
	var count int64 = 0
	var size uint32 = 100
	serviceInfo := make(map[string]model.Instance)
	d.lock.RLock()
	listen := d.listen[groupName]
	d.lock.RUnlock()
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
				serviceInfo[fmt.Sprintf("%s:%d", instance.Ip, instance.Port)] = instance
			}
		}
		count = count + 1
		page++
	}
	for key, value := range serviceInfo {
		if _, ok := listen.Load(key); !ok {
			up(value)
		}
		listen.Store(key, value)
	}
	listen.Range(func(key, value interface{}) bool {
		if _, ok := serviceInfo[key.(string)]; !ok {
			listen.Delete(key)
			down(value.(model.Instance))
		}
		return true
	})
}
