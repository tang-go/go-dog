package param

//SetAdminStatusRequest 设置管理员状态响应
type SetAdminStatusRequest struct {
	//发起请求用户的token
	Token string
	//设置的管理员的ID
	AdminID int64
	//禁用状态
	IsDisable bool
	//路由
	MenuPath string
}

//SetAdminStatusResponse 设置管理员状态响应
type SetAdminStatusResponse struct {
	//结果
	Success bool
}

//GetAdminCodeResponse 获取管理员验证码响响应
type GetAdminCodeResponse struct {
	//ID
	ID string
	//图片
	Img string
}

//LoginAdminRequest 管理员登陆请求
type LoginAdminRequest struct {
	//电话
	Phone string
	//密码
	Pwd string
	//验证码ID
	ID string
	//验证码答案
	Answer string
}

//LoginAdminResponse 管理员登陆响应
type LoginAdminResponse struct {
	//名称
	Name string
	//业主ID
	OwnerID int64
	//注册用户的token
	Token string
}

//RegisterAdminRequest 注册管理员请求
type RegisterAdminRequest struct {
	//电话
	Phone string
	//密码
	Pwd string
}

//RegisterAdminResponse 注册管理员响应
type RegisterAdminResponse struct {
	//注册用户的token
	Success bool
}

//AddAdminRequest 添加管理员账号请求
type AddAdminRequest struct {
	//发起请求用户的token
	Token string
	//电话
	Phone string
	//密码
	Pwd string
	//姓名
	Name string
	//等级
	Level int
	//业主ID
	OwnerID int64
	//路由
	MenuPath string
}

//AddAdminResponse 添加管理员账号响应
type AddAdminResponse struct {
	//成功true 失败false
	Success bool
}

//UpdateAdminRequest 修改管理员账号请求
type UpdateAdminRequest struct {
	//发起请求用户的token
	Token string
	//管理员ID
	AdminID int64
	//电话
	Phone string
	//姓名
	Name string
	//路由
	MenuPath string
}

//UpdateAdminResponse 修改管理员账号响应
type UpdateAdminResponse struct {
	//成功true 失败false
	Success bool
}

//UpdateAdminPwdRequest 修改管理员账号请求
type UpdateAdminPwdRequest struct {
	//发起请求用户的token
	Token string
	//管理员ID
	AdminID int64
	//密码
	Pwd string
	//路由
	MenuPath string
}

//UpdateAdminPwdResponse 修改管理员账号响应
type UpdateAdminPwdResponse struct {
	//成功true 失败false
	Success bool
}

//DelAdminRequest 添加管理员账号请求
type DelAdminRequest struct {
	//发起请求用户的token
	Token string
	//删除管理员ID
	AdminID int64
	//路由
	MenuPath string
}

//DelAdminResponse 添加管理员账号响应
type DelAdminResponse struct {
	//成功true 失败false
	Success bool
}

//GetAdminListRequest 获取管理员账号请求
type GetAdminListRequest struct {
	//发起请求用户的token
	Token string
	//路由
	MenuPath string
}

//GetAdminListResponse 获取管理员列表响应
type GetAdminListResponse struct {
	//消息体
	AdminInfos []*AdminInfo
}

//AdminInfo 管理员信息
type AdminInfo struct {
	//账号自动生成--规程ID加随机6位数
	AdminID int64
	//名称
	Name string
	//电话
	Phone string
	//禁用状态
	IsDisable bool
	//等级
	Level int
	//注册时间
	Time string
}

//GetAuthRequest 获取管理员权限
type GetAuthRequest struct {
	//发起请求用户的token
	Token string
	//管理员ID
	AdminID int64
	//路由
	MenuPath string
}

//GetAuthResponse 获取用户权限
type GetAuthResponse struct {
	//IDs
	MenuIDs []int64
}

//SetAdminAuthRequest 添加管理员菜单菜单请求
type SetAdminAuthRequest struct {
	//发起请求用户的token
	Token string
	//菜单ID
	MenuIDs []int64
	//名称
	AdminID int64
	//路由
	MenuPath string
}

//SetAdminAuthResponse 添加管理员菜单菜单响应
type SetAdminAuthResponse struct {
	//成功true 失败false
	Suceess bool
}
