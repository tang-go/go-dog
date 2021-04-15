package main

import (
	"github.com/tang-go/go-dog/example/dao/services/service-two/router"
	service "github.com/tang-go/go-dog/example/dao/services/service-two/servcie"
	"github.com/tang-go/go-dog/log"
)

func main() {
	s := service.NewService(router.ExampleRouter)
	if e := s.Run(); e != nil {
		log.Errorln(e.Error())
	}
}
