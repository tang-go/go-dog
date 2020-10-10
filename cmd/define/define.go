//Package define 定义框架常量
package define

//服务定义
const (
	// SvcPrefix 统一前缀
	SvcPrefix = "go-dog/"
)

//定义服务
const (
	// SvcAdmin 控制中心
	SvcController = SvcPrefix + "controller"
	// SvcGateWay 网管服务
	SvcGateWay = SvcPrefix + "gateway"
)

//Aes key加密
const (
	AdminAes = "aes.shoplive.code"
)

const (
	//AdminTokenValidityTime 管理员token有效时间 1小时
	AdminTokenValidityTime = 60 * 60
	//CodeValidityTime 验证码有效时间 5分钟
	CodeValidityTime = 60 * 5
)

const (
	//MaxClientRequestCount 客户端最大访问量
	MaxClientRequestCount = 1000000
	//MaxServiceRequestCount 服务端最大访问量
	MaxServiceRequestCount = 1000000
	//TTL 心跳时间
	TTL int64 = 5
)
