package param

//GetCodeReq 获取验证码请求
type GetCodeReq struct {
	Code int64 `description:"随机码" type:"int64"`
}

//GetCodeRes 获取验证码响应
type GetCodeRes struct {
	ID  string `description:"验证码ID" type:"string"`
	Img string `description:"验证码图片" type:"string"`
}
