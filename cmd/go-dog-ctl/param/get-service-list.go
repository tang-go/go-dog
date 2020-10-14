package param

//GetServiceReq 获取服务列表请求
type GetServiceReq struct {
	Code string `description:"业务随机码" type:"string"`
}

//GetServiceRes 获取服务响应
type GetServiceRes struct {
	List []*ServiceInfo `description:"用户token" type:"[]*ServiceInfo"`
}

//ServiceInfo 服务信息
type ServiceInfo struct {
	Key       string    `description:"注册时候使用的唯一key" type:"string"`
	Name      string    `description:"服务名称" type:"string"`
	Address   string    `description:"服务地址" type:"string"`
	Port      int       `description:"端口" type:"int"`
	Methods   []*Method `description:"服务方法" type:"[]*Method"`
	Explain   string    `description:"服务description" type:"string"`
	Longitude int64     `description:"经度" type:"int64"`
	Latitude  int64     `description:"纬度" type:"int64"`
	Time      string    `description:"服务上线时间" type:"string"`
}

//Method 方法
type Method struct {
	Name     string                 `description:"方法名称" type:"string"`
	Level    int8                   `description:"方法等级" type:"int8"`
	Request  map[string]interface{} `description:"请求json格式展示" type:"string"`
	Response map[string]interface{} `description:"响应格式" type:"string"`
	Explain  string                 `description:"方法description" type:"string"`
	IsAuth   bool                   `description:"是否验证" type:"bool"`
}
