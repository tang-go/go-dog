package serviceinfo

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

//APIServiceInfo 提供HTTP API请求的服务（组成路由格式:name/version/path）
type APIServiceInfo struct {
	Key       string //注册时候使用的唯一key
	Name      string //服务名称
	Address   string //服务地址
	Port      int    //端口
	API       []*API //服务方法
	Explain   string //服务说明
	Longitude int64  //经度
	Latitude  int64  //纬度
	Time      string //服务上线时间
}

//API 服务提供的API接口
type API struct {
	Name     string //方法名称
	Level    int8   //方法等级
	Request  string //请求json格式展示
	Response string //响应格式
	Explain  string //方法说明
	IsAuth   bool   //是否验证
	Version  string //版本 例如:v1 v2
	Path     string //http请求路径
	Kind     string //请求类型 POST GET DELETE PUT
}

//Flusing 熔断
type Flusing struct {
	ServiceKey string
	Method     string
}
