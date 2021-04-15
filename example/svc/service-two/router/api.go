package router

import (
	service "github.com/tang-go/go-dog/example/svc/service-two/servcie"
	"github.com/tang-go/go-dog/plugins"
)

//ExampleRouter 演示路由
func ExampleRouter(router plugins.Service, s *service.Service) {
	router.RPC("Add", 4, false, "加法", s.Add)
}
