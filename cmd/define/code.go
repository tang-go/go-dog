package define

//返回码定义
const (
	//SuccessCode 成功返回码
	SuccessCode = 10000
	//GetCodeErr 获取验证码失败
	GetCodeErr = 10001
	//VerfiyCodeErr 验证码失败
	VerfiyCodeErr = 10002
	//AdminLoginErr 管理员登录失败
	AdminLoginErr = 10003
	//AdminTokenErr 管理员Token验证失败
	AdminTokenErr = 10004
)
