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
		shopNoAuth := shopGate.APINoAuth()
		{
			shopV1 := shopNoAuth.APIGroup("演示相关").APIVersion("v1")
			{
				shopV1.APILevel(4).POST("Add", "add", "加法", s.Add)
			}
		}
	}
}
