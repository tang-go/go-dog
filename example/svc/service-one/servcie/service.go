package service

import (
	"github.com/tang-go/go-dog/example/define"
	"github.com/tang-go/go-dog/example/svc/service-one/param"
	twoParam "github.com/tang-go/go-dog/example/svc/service-two/param"
	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/pkg/service"
	"github.com/tang-go/go-dog/plugins"
)

//Service 用户服务
type Service struct {
	service plugins.Service
}

//NewService 初始化服务
func NewService(routers ...func(plugins.Service, *Service)) *Service {
	s := new(Service)
	//初始化服务端
	s.service = service.CreateService(define.ServiceOne)
	//初始化路由
	for _, router := range routers {
		router(s.service, s)
	}
	return s
}

//Run 启动
func (s *Service) Run() error {
	err := s.service.Run()
	if err != nil {
		log.Errorln(err.Error())
	}
	return err
}

//Add 加法计算
func (s *Service) Add(ctx plugins.Context, request param.AddReq) (response param.AddRsp, err error) {
	rsp := new(twoParam.AddRsp)
	if err = ctx.GetClient().Call(ctx, plugins.RandomMode, define.ServiceTwo, "Add", &twoParam.AddReq{
		X: request.X,
		Y: request.Y,
	}, rsp); err != nil {
		log.Errorln(err.Error())
	}
	response.Z = rsp.Z
	log.Traceln("调用ADD")
	return
}
