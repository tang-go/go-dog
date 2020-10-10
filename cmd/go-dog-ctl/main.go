package main

import (
	"go-dog/cmd/go-dog-ctl/api"
)

func main() {
	s := api.NewService()
	s.Run()
}
