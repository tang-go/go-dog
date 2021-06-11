package serviceinfo

//ServiceInfo 服务信息
type ServiceInfo struct {
	Key       string    //唯一组件
	Group     string    //组
	Name      string    //服务名称
	Address   string    //服务地址
	Port      int       //端口
	Explain   string    //服务说明
	Longitude int64     //经度
	Latitude  int64     //纬度
	API       []*API    //api信息
	Methods   []*Method //方法信息
	Time      string    //服务上线时间
}

//Method 方法
type Method struct {
	Name     string                 //方法名称
	Level    int8                   //方法等级
	Request  map[string]interface{} //请求json格式展示
	Response map[string]interface{} //响应格式
	Explain  string                 //方法说明
	IsAuth   bool                   //是否验证
}

//API 服务提供的API接口
type API struct {
	Gate     string                 //注册网关的名称
	Name     string                 //方法名称
	Group    string                 //api的分组
	Level    int8                   //方法等级
	Request  map[string]interface{} //请求json格式展示
	Response map[string]interface{} //响应格式
	Explain  string                 //方法说明
	IsAuth   bool                   //是否验证
	Version  string                 //版本 例如:v1 v2
	Path     string                 //http请求路径
	Kind     string                 //请求类型 POST GET DELETE PUT
}

//Flusing 熔断
type Flusing struct {
	ServiceKey string
	Method     string
}

//ServcieAPI api列表
type ServcieAPI struct {
	Method  *API
	Gate    string //注册网关的名称
	Tags    string
	Name    string
	Explain string
	Count   int32
}
