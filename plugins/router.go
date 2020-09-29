package plugins

import (
	"go-dog/header"
)

//Router RPC路由
type Router interface {
	//RegisterByMethod 注册方法
	RegisterByMethod(name string, fn interface{}) (arg interface{}, reply interface{})

	//Call 调用方法
	Call(ctx Context, argv *header.Request) ([]byte, error)
}
