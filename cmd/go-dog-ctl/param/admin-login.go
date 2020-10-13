package param

//AdminLoginReq 管理员登录
type AdminLoginReq struct {
	Phone  string  `description:"电话" type:"string"`
	Pwd    string  `description:"密码" type:"string"`
	ID     string  `description:"验证码ID" type:"string"`
	Answer string  `description:"验证码答案" type:"string"`
	Test   []int64 `description:"测试" type:"int64"`
}

//AdminLoginRes 管理员登录返回
type AdminLoginRes struct {
	Name    string `description:"名称" type:"string"`
	OwnerID int64  `description:"业主ID" type:"int64"`
	Token   string `description:"注册用户的token" type:"string"`
}