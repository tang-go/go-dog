package param

//SetOwnerStatusRequest 设置业主状态响应
type SetOwnerStatusRequest struct {
	//发起请求用户的token
	Token string
	//设置的业主的ID
	OwnerID int64
	//禁用状态
	IsDisable bool
	//路由
	MenuPath string
}

//SetOwnerStatusResponse 设置业主状态响应
type SetOwnerStatusResponse struct {
	//结果
	Success bool
}

//SetOwnerAuthRequest 添加业主菜单菜单请求
type SetOwnerAuthRequest struct {
	//发起请求用户的token
	Token string
	//菜单ID
	MenuIDs []int64
	//业主
	OwnerID int64
	//路由
	MenuPath string
}

//SetOwnerAuthResponse 添加业主菜单菜单响应
type SetOwnerAuthResponse struct {
	//成功true 失败false
	Suceess bool
}

//GetOwnerAuthRequest 获取业主权限
type GetOwnerAuthRequest struct {
	//发起请求用户的token
	Token string
	//业主ID
	OwnerID int64
	//路由
	MenuPath string
}

//GetOwnerAuthResponse 获取业主权限响应
type GetOwnerAuthResponse struct {
	//IDs
	MenuIDs []int64
}

//AddOwnerRequest 添加业主账号请求
type AddOwnerRequest struct {
	//发起请求用户的token
	Token string
	//电话
	Phone string
	//密码
	Pwd string
	//姓名
	Name string
	//菜单ID
	MenuIDs []int64
	//路由
	MenuPath string
}

//AddOwnerResponse 添加业主账号响应
type AddOwnerResponse struct {
	//成功true 失败false
	Success bool
}

//GetOwnerListRequest 获取业主账号请求
type GetOwnerListRequest struct {
	//发起请求用户的token
	Token string
	//路由
	MenuPath string
}

//GetOwnerListResponse 获取业主列表响应
type GetOwnerListResponse struct {
	//消息体
	OwnerInfo []*OwnerInfo
}

//OwnerInfo 业主信息
type OwnerInfo struct {
	//账号自动生成--规程ID加随机6位数
	OwnerID int64
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

//UpdateOwnerRequest 修改业主账号请求
type UpdateOwnerRequest struct {
	//发起请求用户的token
	Token string
	//业主ID
	OwnerID int64
	//电话
	Phone string
	//姓名
	Name string
	//路由
	MenuPath string
}

//UpdateOwnerResponse 修改业主账号响应
type UpdateOwnerResponse struct {
	//成功true 失败false
	Success bool
}

//UpdateOwnerPwdRequest 修改业主账号请求
type UpdateOwnerPwdRequest struct {
	//发起请求用户的token
	Token string
	//业主ID
	OwnerID int64
	//密码
	Pwd string
	//路由
	MenuPath string
}

//UpdateOwnerPwdResponse 修改业主账号响应
type UpdateOwnerPwdResponse struct {
	//成功true 失败false
	Success bool
}

//DelOwnerRequest 添加业主账号请求
type DelOwnerRequest struct {
	//发起请求用户的token
	Token string
	//业主ID
	OwnerID int64
	//路由
	MenuPath string
}

//DelOwnerResponse 添加业主账号响应
type DelOwnerResponse struct {
	//成功true 失败false
	Success bool
}

//VerifyOwnerRequest 验证业主
type VerifyOwnerRequest struct {
	//业主ID
	OwnerID int64
}

//VerifyOwnerResponse 验证业主响应
type VerifyOwnerResponse struct {
	//成功true 失败false
	Success bool
}
