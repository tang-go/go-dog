package main

import (
	"github.com/tang-go/go-dog/example/define"
	"github.com/tang-go/go-dog/pkg/gateway"
)

// @Tags 权限管理
// @Summary 添加功能1
// @Param  token header string true "token"
// @Router /authority/func/add [post]
func a() {

}

// @Tags 权限管理
// @Summary 添加功能2
// @Param  token header string true "token"
// @Router /authority/func/add [get]
func b() {

}

// @Tags 权限管理
// @Summary 添加功能3
// @Param  token header string true "token"
// @Router /authority/func/add [put]
func c() {

}

// @Tags 权限管理
// @Summary 添加功能4
// @Param  token header string true "token"
// @Router /authority/func/add [delete]
func d() {

}

func main() {
	gate := gateway.NewGateway(define.ExampleGate)
	gate.Run(8081)
}
