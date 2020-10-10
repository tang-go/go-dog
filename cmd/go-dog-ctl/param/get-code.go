package param

//GetCodeReq 获取验证码请求
type GetCodeReq struct {
	Code string `说明:"随机码" 类型:"string"`
}

//GetCodeRes 获取验证码响应
type GetCodeRes struct {
	ID  string `说明:"验证码ID" 类型:"string"`
	Img string `说明:"验证码图片" 类型:"string"`
}
