package param

//GetServiceReq 获取服务列表请求
type GetServiceReq struct {
	Token string
}

//GetServiceRes 获取服务响应
type GetServiceRes struct {
	List []*ServiceInfo
}

//ServiceInfo 服务信息
type ServiceInfo struct {
	Key       string    //注册时候使用的唯一key
	Name      string    //服务名称
	Address   string    //服务地址
	Port      int       //端口
	Methods   []*Method //服务方法
	Explain   string    //服务说明
	Longitude int64     //经度
	Latitude  int64     //纬度
	Time      string    //服务上线时间
}

//Method 方法
type Method struct {
	Name     string //方法名称
	Level    int8   //方法等级
	Request  string //请求json格式展示
	Response string //响应格式
	Explain  string //方法说明
	IsAuth   bool   //是否验证
}
