package consul

import (
	"sync"

	"github.com/hashicorp/consul/api"
)

//Init 初始化
func Init(address string) {
	gOnce.Do(func() {
		gConsul = newConsul(address)
	})
}

//GetRegister 获取注册中心
func GetRegister() *Register {
	return gConsul.register
}

//GetDiscovery 获取服务发现
func GetDiscovery() *Discovery {
	return gConsul.discovery
}

var (
	gConsul *Consul
	gOnce   sync.Once
)

type Consul struct {
	register  *Register
	discovery *Discovery
}
type Address struct {
	IP   string
	Port uint64
}

//newConsul 初始化私有consul对象
func newConsul(address string) *Consul {
	// 创建连接consul服务配置
	config := api.DefaultConfig()
	config.Address = address
	client, err := api.NewClient(config)
	if err != nil {
		panic(err.Error())
	}
	c := new(Consul)
	c.discovery = newDiscovery(client)
	c.register = newRegister(client)
	return c
}
