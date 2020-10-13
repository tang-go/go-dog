package param

//GetCodeReq 获取验证码请求
type GetCodeReq struct {
	Code string `description:"随机码" type:"string"`
}

//GetCodeRes 获取验证码响应
type GetCodeRes struct {
	ID  string `description:"验证码ID" type:"string"`
	Img string `description:"验证码图片" type:"string"`
}
