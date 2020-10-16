package customerror

import "time"

const (
	//SuccessCode 成功
	SuccessCode = 200
	//ConnectClose 链接关闭
	ConnectClose = 400
	//RPCNotFind 没有找到方法
	RPCNotFind = 404
	//RequestTimeout 请求超时
	RequestTimeout = 408
	//InternalServerError 服务错误
	InternalServerError = 500
	//UnknownError 未知错误
	UnknownError = 505
	//ClientLimitError 客户端限流
	ClientLimitError = 506
	//SeviceLimitError 服务端限流
	SeviceLimitError = 507
	//ParamError 参数错误
	ParamError = 508
)

//Error 错误定义
type Error struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Time int64  `json:"time"`
}

//EnCodeError 创建一个错误
func EnCodeError(code int, msg string) *Error {
	return &Error{
		Code: code,
		Msg:  msg,
		Time: time.Now().Unix(),
	}
}

//DeCodeError 解析错误
func DeCodeError(e error) *Error {
	if e == nil {
		return nil
	}
	if err, ok := e.(*Error); ok {
		return err
	}
	return EnCodeError(UnknownError, e.Error())
}

//Error 输出错误
func (e *Error) Error() string {
	return e.Msg
}
