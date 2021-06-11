package plugins

import "github.com/tang-go/go-dog/serviceinfo"

//Selector 路由选择器
type Selector interface {

	//GetByAddress 通过地址获取rpc服务信息
	GetByAddress(discovery Discovery, address string, fusing Fusing, name string, method string) (*serviceinfo.ServiceInfo, error)

	//RandomMode 随机模式(失败即返回)
	RandomMode(discovery Discovery, fusing Fusing, name string, method string) (*serviceinfo.ServiceInfo, error)

	//RangeMode 遍历模式(一个返回成功,或者全部返回失败)
	RangeMode(discovery Discovery, fusing Fusing, name string, method string, f func(*serviceinfo.ServiceInfo) bool) error

	//HashMode 通过hash值访问一个服务(失败即返回)
	HashMode(discovery Discovery, fusing Fusing, name string, method string) (*serviceinfo.ServiceInfo, error)

	//Custom 自定义 --目前默认随机
	Custom(discovery Discovery, fusing Fusing, name string, method string) (*serviceinfo.ServiceInfo, error)
}
