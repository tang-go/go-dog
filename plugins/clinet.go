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

	//GetCodec 获取编码插件
	GetCodec() Codec

	//GetCfg 获取配置
	GetCfg() Cfg

	//GetDiscovery 获取服务发现
	GetDiscovery() Discovery

	//GetFusing 获取熔断插件
	GetFusing() Fusing

	//GetLimit 获取限流插件
	GetLimit() Limit

	//GetAllService 获取所有服务
	GetAllService() (services []*serviceinfo.ServiceInfo)

	//Call 调用函数
	Call(ctx Context, mode Mode, name string, method string, args interface{}, reply interface{}) error

	//SendRequest 发生请求
	SendRequest(ctx Context, mode Mode, name string, method string, code string, args []byte) (reply []byte, e error)

	//Close 关闭
	Close()
}
