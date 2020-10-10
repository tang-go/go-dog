package param

//GetCodeReq 获取验证码请求
type GetCodeReq struct {
}

//GetCodeRes 获取验证码响应
type GetCodeRes struct {
	//ID
	ID string
	//图片
	Img string
}
