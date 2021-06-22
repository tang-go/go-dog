package plugins

import (
	"github.com/tang-go/go-dog/serviceinfo"
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

	//GetAllRPCService 获取所有RPC服务
	GetAllRPCService() (services []*serviceinfo.ServiceInfo)

	//GetAllAPIService 获取所有API服务
	GetAllAPIService() (services []*serviceinfo.ServiceInfo)

	//Call 调用函数
	Call(ctx Context, mode Mode, name string, method string, args interface{}, reply interface{}) error

	//Broadcast 广播
	Broadcast(ctx Context, name string, method string, args interface{}, reply interface{}) error

	//SendRequest 发生请求
	SendRequest(ctx Context, mode Mode, name string, method string, code string, args []byte) (reply []byte, e error)

	//CallByAddress 指定地址调用
	CallByAddress(ctx Context, address string, name string, method string, args interface{}, reply interface{}) error

	//Close 关闭
	Close()
}
