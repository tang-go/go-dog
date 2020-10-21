package plugins

import "github.com/tang-go/go-dog/serviceinfo"

//Selector 路由选择器
type Selector interface {

	//RandomMode 随机模式(失败即返回)
	RandomMode(discovery Discovery, fusing Fusing, name string, method string) (*serviceinfo.RPCServiceInfo, error)

	//RangeMode 遍历模式(一个返回成功,或者全部返回失败)
	RangeMode(discovery Discovery, fusing Fusing, name string, method string, f func(*serviceinfo.RPCServiceInfo) bool) error

	//HashMode 通过hash值访问一个服务(失败即返回)
	HashMode(discovery Discovery, fusing Fusing, name string, method string) (*serviceinfo.RPCServiceInfo, error)

	//Custom 自定义 --目前默认随机
	Custom(discovery Discovery, fusing Fusing, name string, method string) (*serviceinfo.RPCServiceInfo, error)
}
