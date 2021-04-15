package client

import (
	"sync"

	"github.com/tang-go/go-dog/example/dao/services/service-two/client"
)

var (
	only *Client
	once sync.Once
)

//Client rpc接口
type Client struct {
	serviceTwo *client.ServiceTwo
}

func newClient() *Client {
	rpc := new(Client)
	rpc.serviceTwo = client.NewServiceTwo()
	return rpc
}

//GetServiceTwo 获取对象
func (c *Client) GetServiceTwo() *client.ServiceTwo {
	return c.serviceTwo
}

// Only 初始化服务
func Only() *Client {
	once.Do(func() {
		only = new(Client)
	})
	return only
}
