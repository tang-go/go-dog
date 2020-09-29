package main

import (
	"go-dog/cmd/go-dog-ctl/controller"
)

func main() {
	s := controller.NewController()
	s.Run()
}
