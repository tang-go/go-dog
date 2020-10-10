package param

//GetServiceReq 获取服务列表请求
type GetServiceReq struct {
	Token string `说明:"用户token" 类型:"string"`
}

//GetServiceRes 获取服务响应
type GetServiceRes struct {
	List []*ServiceInfo `说明:"用户token" 类型:"[]*ServiceInfo"`
}

//ServiceInfo 服务信息
type ServiceInfo struct {
	Key       string    `说明:"注册时候使用的唯一key" 类型:"string"`
	Name      string    `说明:"服务名称" 类型:"string"`
	Address   string    `说明:"服务地址" 类型:"string"`
	Port      int       `说明:"端口" 类型:"int"`
	Methods   []*Method `说明:"服务方法" 类型:"[]*Method"`
	Explain   string    `说明:"服务说明" 类型:"string"`
	Longitude int64     `说明:"经度" 类型:"int64"`
	Latitude  int64     `说明:"纬度" 类型:"int64"`
	Time      string    `说明:"服务上线时间" 类型:"string"`
}

//Method 方法
type Method struct {
	Name     string                 `说明:"方法名称" 类型:"string"`
	Level    int8                   `说明:"方法等级" 类型:"int8"`
	Request  map[string]interface{} `说明:"请求json格式展示" 类型:"string"`
	Response map[string]interface{} `说明:"响应格式" 类型:"string"`
	Explain  string                 `说明:"方法说明" 类型:"string"`
	IsAuth   bool                   `说明:"是否验证" 类型:"bool"`
}
