package router

import (
	service "github.com/tang-go/go-dog/example/dao/services/service-two/servcie"
	"github.com/tang-go/go-dog/plugins"
)

//ExampleRouter 演示路由
func ExampleRouter(router plugins.Service, s *service.Service) {
	router.RPC().Level(4).NoAuth().Method("Add", "加法", s.Add)
}
