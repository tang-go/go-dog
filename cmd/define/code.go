package define

//返回码定义
const (
	//SuccessCode 成功返回码
	SuccessCode = 10000
	//SystemErr 系统错误
	SystemErr = 10001
	//TokenErr token失败
	TokenErr = 10002
	//TokenOver token过期
	TokenOver = 10003
	//PermissionErr 权限错误
	PermissionErr = 10004
	//GetCodeErr 获取验证码失败
	GetCodeErr = 10005
	//VerifyCodeErr 验证验证码失败
	VerifyCodeErr = 10006
	//PhoneErr 手机号不正确
	PhoneErr = 10007
	//PhoneIsNot 手机号不存在
	PhoneIsNot = 10008
	//PwdErr 密码错误
	PwdErr = 10009
	//AddAdminErr 添加管理员错误
	AddAdminErr = 10010
	//ParamErr 参数错误
	ParamErr = 10011
)
