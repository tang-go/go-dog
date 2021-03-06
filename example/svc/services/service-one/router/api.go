package router

import (
	"github.com/tang-go/go-dog/example/define"
	service "github.com/tang-go/go-dog/example/svc/services/service-one/servcie"
	"github.com/tang-go/go-dog/plugins"
)

//ExampleRouter 演示路由
func ExampleRouter(router plugins.Service, s *service.Service) {
	shopGate := router.HTTP(define.ExampleGate)
	{
		shopNoAuth := shopGate.NoAuth()
		{
			shopV1 := shopNoAuth.Group("演示相关").Version("v1")
			{
				shopV1.Level(4).POST("Add", "add", "加法", s.Add)
				shopV1.Level(4).PUT("Add", "add", "加法", s.Add)
				shopV1.Level(4).GET("Add", "add", "加法", s.Add)
				shopV1.Level(4).DELETE("Add", "add", "加法", s.Add)
			}
		}
	}
}
