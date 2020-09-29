package param

//GetAPIListReq 获取API列表请求
type GetAPIListReq struct {
	Token string
}

//GetAPIListRes 获取列表返回
type GetAPIListRes struct {
	List []*Service
}

//Service 服务
type Service struct {
	APIS    []*API
	Name    string
	Explain string
}

//API 服务提供的API接口
type API struct {
	Name     string //服务名称
	Level    int8   //方法等级
	Request  string //请求json格式展示
	Response string //响应格式
	Explain  string //方法说明
	IsAuth   bool   //是否验证
	Version  string //版本 例如:v1 v2
	URL      string //http请求路径
	Kind     string //请求类型 POST GET DELETE PUT
}
