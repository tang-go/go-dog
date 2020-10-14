package api

//InitAPI 初始化API
func (pointer *Service) InitAPI() {
	//获取图片验证码
	pointer.service.GET("GetCode", "v1", "get/code",
		3,
		false,
		"获取图片验证码",
		pointer.GetCode)
	//验证码验证码
	pointer.service.POST("AdminLogin", "v1", "admin/login",
		3,
		false,
		"管理员登录",
		pointer.AdminLogin)
	//获取API列表
	pointer.service.GET("GetAPIList", "v1", "get/api/list",
		3,
		true,
		"获取api列表",
		pointer.GetAPIList)
	//获取服务列表
	pointer.service.GET("GetServiceList", "v1", "get/service/list",
		3,
		true,
		"获取服务列表",
		pointer.GetServiceList)
}
