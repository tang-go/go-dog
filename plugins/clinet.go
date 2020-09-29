package plugins

import (
	"go-dog/serviceinfo"
)

//Mode 模式
type Mode int8

const (
	//RandomMode 随机模式
	RandomMode Mode = iota
	//RangeMode 遍历模式
	RangeMode
	//HashMode 一致性hash模式
	HashMode
)

//Client 客户端
type Client interface {

	//GetAllService 获取所有服务
	GetAllService() (services []*serviceinfo.ServiceInfo)

	//SetFlowLimit 设置最大流量限制
	SetFlowLimit(max int64)

	//ServiceOnlineNotice 服务上线
	ServiceOnlineNotice(key string, info *serviceinfo.ServiceInfo)

	//ServiceOfflineNotice 服务下线
	ServiceOfflineNotice(key string)

	//Call 调用函数
	Call(ctx Context, mode Mode, name string, method string, args interface{}, reply interface{}) error

	//SendRequest 发生请求
	SendRequest(ctx Context, mode Mode, name string, method string, args []byte) (reply []byte, e error)

	//Close 关闭
	Close()
}
