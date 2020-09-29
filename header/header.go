package header

import (
	"go-dog/error"
)

//Request MsgPack请求
type Request struct {
	MsgPackCode
	TTL     int64
	TimeOut int64
	TraceID string
	IsTest  bool
	Address string
	Data    map[string]interface{}
	ID      string
	Name    string
	Method  string
	Arg     []byte
}

//Response MsgPack响应
type Response struct {
	MsgPackCode
	ID     string
	Name   string
	Method string
	Reply  []byte
	Error  *customerror.Error
}

type key int

const (
	//ContextTraceIDValue context TraceID链路追踪ID
	ContextTraceIDValue key = iota
	//ContextIsTestValue 是否为测试数据
	ContextIsTestValue
	//ContextAddressValue 请求客户端IP的值
	ContextAddressValue
	//ContextDataValue context 自定义data值
	ContextDataValue
)
