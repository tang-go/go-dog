package param

//AdminLoginReq 管理员登录
type AdminLoginReq struct {
	Phone  string `说明:"电话" 类型:"string"`
	Pwd    string `说明:"密码" 类型:"string"`
	ID     string `说明:"验证码ID" 类型:"string"`
	Answer string `说明:"验证码答案" 类型:"string"`
}

//AdminLoginRes 管理员登录返回
type AdminLoginRes struct {
	Name    string `说明:"名称" 类型:"string"`
	OwnerID int64  `说明:"业主ID" 类型:"int64"`
	Token   string `说明:"注册用户的token" 类型:"string"`
}
