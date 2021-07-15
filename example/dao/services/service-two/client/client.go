package client

import (
	"github.com/tang-go/go-dog/example/dao/services/service-two/param"
	"github.com/tang-go/go-dog/example/define"
	"github.com/tang-go/go-dog/plugins"
)

//ServiceTwo 演示案例二
type ServiceTwo struct {
}

//NewServiceTwo 初始化
func NewServiceTwo() *ServiceTwo {
	s := new(ServiceTwo)
	return s
}

//Add 添加
func (s *ServiceTwo) Add(ctx plugins.Context, x, y int64) (z int64, err error) {
	rsp := new(param.AddRsp)
	err = ctx.GetClient().Call(ctx, plugins.RandomMode, define.ServiceTwo, "", "Add", &param.AddReq{
		X: x,
		Y: y,
	}, rsp)
	return rsp.Z, err
}
