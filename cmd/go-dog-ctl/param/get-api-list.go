package param

//GetAPIListReq 获取API列表请求
type GetAPIListReq struct {
	Code string `json:"code" description:"业务随机码" type:"string"`
}

//GetAPIListRes 获取列表返回
type GetAPIListRes struct {
	List []*Service `json:"list" description:"API服务列表" type:"[]*Service"`
}

//Service 服务
type Service struct {
	APIS    []*API `json:"apis" description:"API集合" type:"[]*API"`
	Name    string `json:"name" description:"API服务名称" type:"string"`
	Explain string `json:"explain" description:"API服务description" type:"string"`
}

//API 服务提供的API接口
type API struct {
	Name     string                 `json:"name" description:"服务名称" type:"string"`
	Level    int8                   `json:"level" description:"方法等级" type:"int8"`
	Request  map[string]interface{} `json:"request" description:"请求json格式展示" type:"string"`
	Response map[string]interface{} `json:"response" description:"响应格式" type:"string"`
	Explain  string                 `json:"explain" description:"方法description" type:"string"`
	IsAuth   bool                   `json:"isAuth" description:"是否验证" type:"bool"`
	Version  string                 `json:"version" description:"版本 例如:v1 v2" type:"string"`
	URL      string                 `json:"url" description:"http请求路径" type:"string"`
	Kind     string                 `json:"kind" description:"请求type POST GET DELETE PUT" type:"string"`
}
