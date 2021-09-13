package consul

import (
	"context"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/tang-go/go-dog/lib/uuid"
	"github.com/tang-go/go-dog/log"
)

//注册参数
type RegisterInstanceParam struct {
	Name              string
	Tags              []string
	Port              int
	Address           string
	EnableTagOverride bool
	Meta              map[string]string
}

//Register 注册中心
type Register struct {
	client        *api.Client
	registrations []*api.AgentServiceRegistration
	ctx           context.Context
	cancel        context.CancelFunc
}

//NewRegister 创建注册中心
func newRegister(c *api.Client) *Register {
	r := new(Register)
	r.client = c
	r.ctx, r.cancel = context.WithCancel(context.Background())
	return r
}

//Register 注册一个服务
func (r *Register) Register(info RegisterInstanceParam) error {
	// 创建注册到consul的服务到
	registration := new(api.AgentServiceRegistration)
	registration.ID = uuid.GetToken()
	registration.Tags = info.Tags
	registration.Name = info.Name
	registration.Port = info.Port
	registration.EnableTagOverride = info.EnableTagOverride
	registration.Address = info.Address
	registration.Meta = info.Meta

	check := new(api.AgentServiceCheck)
	check.CheckID = uuid.GetToken()
	check.TTL = "5s"
	check.Name = info.Name
	check.DeregisterCriticalServiceAfter = "10s"
	registration.Check = check
	// 注册服务到consul
	if err := r.client.Agent().ServiceRegister(registration); err != nil {
		return err
	}
	go func() {
		timeTicker := time.NewTicker(2 * time.Second)
		for {
			select {
			case <-timeTicker.C:
				if err := r.client.Agent().PassTTL(check.CheckID, ""); err != nil {
					log.Errorln(err.Error())
				}
			case <-r.ctx.Done():
				log.Traceln("consul register exit")
				return
			}
		}
	}()
	r.registrations = append(r.registrations, registration)
	return nil
}

//DeregisterInstance 取消注册一个服务
func (r *Register) DeregisterInstance() {
	r.cancel()
	for _, registration := range r.registrations {
		r.client.Agent().ServiceDeregister(registration.ID)
	}
}
