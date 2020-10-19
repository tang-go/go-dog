package main

import "go-dog/cmd/go-dog-find/service"

func main() {
	s := service.NewService()
	s.Run()
}
