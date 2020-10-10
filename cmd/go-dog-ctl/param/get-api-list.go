package param

//GetAPIListReq 获取API列表请求
type GetAPIListReq struct {
	Token string `说明:"用户token" 类型:"string"`
}

//GetAPIListRes 获取列表返回
type GetAPIListRes struct {
	List []*Service `说明:"API服务列表" 类型:"[]*Service"`
}

//Service 服务
type Service struct {
	APIS    []*API `说明:"API集合" 类型:"[]*API"`
	Name    string `说明:"API服务名称" 类型:"string"`
	Explain string `说明:"API服务说明" 类型:"string"`
}

//API 服务提供的API接口
type API struct {
	Name     string                 `说明:"服务名称" 类型:"string"`
	Level    int8                   `说明:"方法等级" 类型:"int8"`
	Request  map[string]interface{} `说明:"请求json格式展示" 类型:"string"`
	Response map[string]interface{} `说明:"响应格式" 类型:"string"`
	Explain  string                 `说明:"方法说明" 类型:"string"`
	IsAuth   bool                   `说明:"是否验证" 类型:"bool"`
	Version  string                 `说明:"版本 例如:v1 v2" 类型:"string"`
	URL      string                 `说明:"http请求路径" 类型:"string"`
	Kind     string                 `说明:"请求类型 POST GET DELETE PUT" 类型:"string"`
}
