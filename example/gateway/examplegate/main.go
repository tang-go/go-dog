package main

import (
	"github.com/tang-go/go-dog/example/define"
	"github.com/tang-go/go-dog/pkg/gateway"
)

func main() {
	gate := gateway.NewGateway(define.ExampleGate)
	gate.Run(8081)
}
