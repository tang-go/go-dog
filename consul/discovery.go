package consul

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/tang-go/go-dog/log"
)

//Instance 实例内容
type Instance struct {
	ID                string
	Service           string
	Tags              []string
	Meta              map[string]string
	Port              int
	Address           string
	EnableTagOverride bool
	CreateIndex       uint64
	ModifyIndex       uint64
	ContentHash       string
}

//Discovery 服务发现
type Discovery struct {
	client *api.Client
}

//NewDiscovery 创建一个服务发现
func newDiscovery(c *api.Client) *Discovery {
	d := new(Discovery)
	d.client = c
	return d
}

//Discovery 服务发现
func (d *Discovery) Discovery(ctx context.Context, tag []string, up func(Instance), down func(Instance)) {
	go func() {
		var listen sync.Map
		ticker := time.NewTicker(time.Second * 2)
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				services, err := d.client.Agent().Services()
				if err != nil {
					log.Errorln(err.Error())
					break
				}
				serviceInfo := make(map[string]api.AgentService)
				for _, service := range services {
					if d.IsContain(service.Tags, tag) == false {
						continue
					}
					result, _, err := d.client.Agent().AgentHealthServiceByID(service.ID)
					if err != nil {
						log.Errorln(err.Error())
						continue
					}
					if result == api.HealthPassing {
						serviceInfo[fmt.Sprintf("%s:%d", service.Address, service.Port)] = *service
					}
				}

				for key, value := range serviceInfo {
					if _, ok := listen.Load(key); !ok {
						up(Instance{
							ID:                value.ID,
							Service:           value.Service,
							Tags:              value.Tags,
							Meta:              value.Meta,
							Port:              value.Port,
							Address:           value.Address,
							EnableTagOverride: value.EnableTagOverride,
							CreateIndex:       value.CreateIndex,
							ModifyIndex:       value.ModifyIndex,
							ContentHash:       value.ContentHash,
						})
					}
					listen.Store(key, value)
				}
				listen.Range(func(key, info interface{}) bool {
					if _, ok := serviceInfo[key.(string)]; !ok {
						listen.Delete(key)
						value := info.(api.AgentService)
						down(Instance{
							ID:                value.ID,
							Service:           value.Service,
							Tags:              value.Tags,
							Meta:              value.Meta,
							Port:              value.Port,
							Address:           value.Address,
							EnableTagOverride: value.EnableTagOverride,
							CreateIndex:       value.CreateIndex,
							ModifyIndex:       value.ModifyIndex,
							ContentHash:       value.ContentHash,
						})
					}
					return true
				})
			}
		}
	}()
}

//IsContain 判断数组是否包含
func (d *Discovery) IsContain(a []string, b []string) bool {
	for _, item := range b {
		ret := false
		for _, i := range a {
			if item == i {
				ret = true
				break
			}
		}
		if ret == false {
			return false
		}
	}
	return true
}
