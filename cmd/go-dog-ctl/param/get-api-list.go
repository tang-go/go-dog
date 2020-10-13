package param

//GetAPIListReq 获取API列表请求
type GetAPIListReq struct {
	Token string `description:"用户token" type:"string"`
}

//GetAPIListRes 获取列表返回
type GetAPIListRes struct {
	List []*Service `description:"API服务列表" type:"[]*Service"`
}

//Service 服务
type Service struct {
	APIS    []*API `description:"API集合" type:"[]*API"`
	Name    string `description:"API服务名称" type:"string"`
	Explain string `description:"API服务description" type:"string"`
}

//API 服务提供的API接口
type API struct {
	Name     string                 `description:"服务名称" type:"string"`
	Level    int8                   `description:"方法等级" type:"int8"`
	Request  map[string]interface{} `description:"请求json格式展示" type:"string"`
	Response map[string]interface{} `description:"响应格式" type:"string"`
	Explain  string                 `description:"方法description" type:"string"`
	IsAuth   bool                   `description:"是否验证" type:"bool"`
	Version  string                 `description:"版本 例如:v1 v2" type:"string"`
	URL      string                 `description:"http请求路径" type:"string"`
	Kind     string                 `description:"请求type POST GET DELETE PUT" type:"string"`
}
