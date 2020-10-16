package param

//GetServiceReq 获取服务列表请求
type GetServiceReq struct {
	Code string `json:"code" description:"业务随机码" type:"string"`
}

//GetServiceRes 获取服务响应
type GetServiceRes struct {
	List []*ServiceInfo `json:"list" description:"用户token" type:"[]*ServiceInfo"`
}

//ServiceInfo 服务信息
type ServiceInfo struct {
	Key       string    `json:"key" description:"注册时候使用的唯一key" type:"string"`
	Name      string    `json:"name" description:"服务名称" type:"string"`
	Address   string    `json:"address" description:"服务地址" type:"string"`
	Port      int       `json:"port" description:"端口" type:"int"`
	Methods   []*Method `json:"methods" description:"服务方法" type:"[]*Method"`
	Explain   string    `json:"explain" description:"服务description" type:"string"`
	Longitude int64     `json:"longitude" description:"经度" type:"int64"`
	Latitude  int64     `json:"latitude" description:"纬度" type:"int64"`
	Time      string    `json:"time" description:"服务上线时间" type:"string"`
}

//Method 方法
type Method struct {
	Name     string                 `json:"name" description:"方法名称" type:"string"`
	Level    int8                   `json:"level" description:"方法等级" type:"int8"`
	Request  map[string]interface{} `json:"request" description:"请求json格式展示" type:"string"`
	Response map[string]interface{} `json:"response" description:"响应格式" type:"string"`
	Explain  string                 `json:"explain" description:"方法description" type:"string"`
	IsAuth   bool                   `json:"isAuth" description:"是否验证" type:"bool"`
}
