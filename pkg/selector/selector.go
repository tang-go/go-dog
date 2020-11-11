package selector

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	customerror "github.com/tang-go/go-dog/error"
	"github.com/tang-go/go-dog/plugins"
	"github.com/tang-go/go-dog/serviceinfo"
)

//Selector 选择器
type Selector struct {
	rnd  *rand.Rand
	lock sync.RWMutex
}

//NewSelector 新建一个选择器
func NewSelector() *Selector {
	s := new(Selector)
	s.rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
	return s
}

//GetByAddress 通过地址获取rpc服务信息
func (s *Selector) GetByAddress(discovery plugins.Discovery, address string, fusing plugins.Fusing, name string, method string) (*serviceinfo.RPCServiceInfo, error) {
	services := discovery.GetRPCServiceByName(name)
	for _, service := range services {
		if !fusing.IsFusing(service.Key, method) {
			if fmt.Sprintf("%s:%d", service.Address, service.Port) == address {
				return service, nil
			}
		}
	}
	return nil, customerror.EnCodeError(customerror.InternalServerError, "没有可用服务")
}

//RandomMode 随机模式(失败即返回)
func (s *Selector) RandomMode(discovery plugins.Discovery, fusing plugins.Fusing, name string, method string) (*serviceinfo.RPCServiceInfo, error) {
	var rpc []*serviceinfo.RPCServiceInfo
	services := discovery.GetRPCServiceByName(name)
	for _, service := range services {
		if !fusing.IsFusing(service.Key, method) {
			rpc = append(rpc, service)
		}
	}
	count := len(rpc)
	if count <= 0 {
		return nil, customerror.EnCodeError(customerror.InternalServerError, "没有可用服务")
	}
	s.lock.Lock()
	index := s.rnd.Intn(count)
	s.lock.Unlock()
	return rpc[index], nil
}

//RangeMode 遍历模式(一个返回成功,或者全部返回失败))
func (s *Selector) RangeMode(discovery plugins.Discovery, fusing plugins.Fusing, name string, method string, f func(*serviceinfo.RPCServiceInfo) bool) error {
	services := discovery.GetRPCServiceByName(name)
	count := len(services)
	if count <= 0 {
		return customerror.EnCodeError(customerror.InternalServerError, "没有可用服务")
	}
	for _, service := range services {
		if !fusing.IsFusing(service.Key, method) {
			if f(service) == true {
				break
			}
		}
	}
	return nil
}

//HashMode 通过hash值访问一个服务(失败即返回)
func (s *Selector) HashMode(discovery plugins.Discovery, fusing plugins.Fusing, name string, method string) (*serviceinfo.RPCServiceInfo, error) {
	return nil, customerror.EnCodeError(customerror.InternalServerError, "暂时没有开启一致性hash")
}

//Custom 自定义 --目前默认随机
func (s *Selector) Custom(discovery plugins.Discovery, fusing plugins.Fusing, name string, method string) (*serviceinfo.RPCServiceInfo, error) {
	return s.RandomMode(discovery, fusing, name, method)
}
