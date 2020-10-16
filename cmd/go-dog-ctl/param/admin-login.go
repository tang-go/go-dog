package param

//AdminLoginReq 管理员登录
type AdminLoginReq struct {
	Phone string `json:"phone" description:"电话" type:"string"`
	Pwd   string `json:"pwd" description:"密码" type:"string"`
}

//AdminLoginRes 管理员登录返回
type AdminLoginRes struct {
	Name    string `json:"name" description:"名称" type:"string"`
	OwnerID int64  `json:"ownerId" description:"业主ID" type:"int64"`
	Token   string `json:"token" description:"注册用户的token" type:"string"`
}
