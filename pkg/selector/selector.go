package selector

import (
	"go-dog/error"
	"go-dog/plugins"
	"go-dog/serviceinfo"
	"math/rand"
	"sync"
	"time"
)

//Selector 选择器
type Selector struct {
	services map[string]*serviceinfo.ServiceInfo
	rnd      *rand.Rand
	lock     sync.RWMutex
}

//NewSelector 新建一个选择器
func NewSelector() *Selector {
	s := new(Selector)
	s.services = make(map[string]*serviceinfo.ServiceInfo)
	s.rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
	return s
}

//GetAllService 获取所有服务
func (s *Selector) GetAllService() (services []*serviceinfo.ServiceInfo) {
	s.lock.RLock()
	for _, service := range s.services {
		services = append(services, service)
	}
	s.lock.RUnlock()
	return
}

//AddService 新增一个服务
func (s *Selector) AddService(key string, service *serviceinfo.ServiceInfo) {
	s.lock.Lock()
	s.services[key] = service
	s.lock.Unlock()
}

//DelService 删除服务
func (s *Selector) DelService(key string) {
	s.lock.Lock()
	delete(s.services, key)
	s.lock.Unlock()
}

//RandomMode 随机模式(失败即返回)
func (s *Selector) RandomMode(fusing plugins.Fusing, name string, method string) (*serviceinfo.ServiceInfo, error) {
	var services []*serviceinfo.ServiceInfo
	s.lock.RLock()
	for _, service := range s.services {
		if service.Name == name {
			if !fusing.IsFusing(service.Key, method) {
				services = append(services, service)
			}
		}
	}
	s.lock.RUnlock()
	count := len(services)
	if count <= 0 {
		return nil, customerror.EnCodeError(customerror.InternalServerError, "没有可用服务")
	}

	s.lock.Lock()
	index := s.rnd.Intn(count)
	s.lock.Unlock()
	return services[index], nil
}

//RangeMode 遍历模式(一个返回成功,或者全部返回失败))
func (s *Selector) RangeMode(fusing plugins.Fusing, name string, method string, f func(*serviceinfo.ServiceInfo) bool) {
	s.lock.RLock()
	for _, service := range s.services {
		if service.Name == name {
			if !fusing.IsFusing(service.Key, method) {
				if f(service) == true {
					break
				}
			}
		}
	}
	s.lock.RUnlock()
}

//HashMode 通过hash值访问一个服务(失败即返回)
func (s *Selector) HashMode(fusing plugins.Fusing, name string, method string) (*serviceinfo.ServiceInfo, error) {
	return nil, customerror.EnCodeError(customerror.InternalServerError, "暂时没有开启一致性hash")
}

//Custom 自定义 --目前默认随机
func (s *Selector) Custom(fusing plugins.Fusing, name string, method string) (*serviceinfo.ServiceInfo, error) {
	return s.RandomMode(fusing, name, method)
}
