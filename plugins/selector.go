package plugins

import "go-dog/serviceinfo"

//Selector 路由选择器
type Selector interface {

	//GetAllService 获取所有服务
	GetAllService() (services []*serviceinfo.ServiceInfo)

	//AddService 新增一个服务
	AddService(key string, service *serviceinfo.ServiceInfo)

	//DelService 删除服务
	DelService(key string)

	//RandomMode 随机模式(失败即返回)
	RandomMode(fusing Fusing, name string, method string) (*serviceinfo.ServiceInfo, error)

	//RangeMode 遍历模式(一个返回成功,或者全部返回失败)
	RangeMode(fusing Fusing, name string, method string, f func(*serviceinfo.ServiceInfo) bool)

	//HashMode 通过hash值访问一个服务(失败即返回)
	HashMode(fusing Fusing, name string, method string) (*serviceinfo.ServiceInfo, error)

	//Custom 自定义 --目前默认随机
	Custom(fusing Fusing, name string, method string) (*serviceinfo.ServiceInfo, error)
}
