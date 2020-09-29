package main

import "go-dog/cmd/go-dog-gw/gateway"

func main() {
	gate := gateway.NewGateway()
	gate.Run()
}
