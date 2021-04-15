package main

import (
	"github.com/tang-go/go-dog/example/svc/service-one/router"
	service "github.com/tang-go/go-dog/example/svc/service-one/servcie"
	"github.com/tang-go/go-dog/log"
)

func main() {
	s := service.NewService(router.ExampleRouter)
	if e := s.Run(); e != nil {
		log.Errorln(e.Error())
	}
}
