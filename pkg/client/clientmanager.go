package client

import (
	"fmt"
	"net"
	"sync"

	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/pkg/rpc"
	"github.com/tang-go/go-dog/plugins"
	"github.com/tang-go/go-dog/serviceinfo"
)

type errService struct {
	count int32
	tm    int64
}

//ManagerClient 管理
type ManagerClient struct {
	codec   plugins.Codec
	clients map[string]*rpc.ClientRPC
	lock    sync.RWMutex
}

//NewManagerClient 创建manager
func NewManagerClient(codec plugins.Codec) *ManagerClient {
	m := new(ManagerClient)
	m.clients = make(map[string]*rpc.ClientRPC)
	m.codec = codec
	return m
}

//GetClient 获取客户端
func (m *ManagerClient) GetClient(service *serviceinfo.RPCServiceInfo) (*rpc.ClientRPC, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	client, ok := m.clients[service.Key]
	if !ok {
		tcpAddr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d", service.Address, service.Port))
		if err != nil {
			log.Errorln(err.Error())
			return nil, err
		}
		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			log.Errorln(err.Error())
			return nil, err
		}
		//创建一个新的链接
		cli := rpc.NewClientRPC(conn, m.codec, func(net.Conn) {
			m.DelClient(service.Key)
			log.Traceln("服务器重启或者崩溃", service.Key)
		})
		m.clients[service.Key] = cli
		return cli, nil
	}
	return client, nil
}

//DelClient 删除客户端
func (m *ManagerClient) DelClient(key string) {
	m.lock.Lock()
	client, ok := m.clients[key]
	if ok {
		client.Close()
		delete(m.clients, key)
	}
	m.lock.Unlock()
}

//Close 关闭
func (m *ManagerClient) Close() {
	m.lock.RLock()
	for _, client := range m.clients {
		client.Close()
	}
	m.lock.RUnlock()
}
