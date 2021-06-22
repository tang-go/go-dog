package plugins

import "github.com/tang-go/go-dog/pkg/config"

const (
	//本地配置模式
	LocalModel = "local"
	//nacos配置模式
	NacosModel = "nacos"
)

//Cfg 配置接口
type Cfg interface {
	//GetServerName 获取服务名称
	//GetServerName() string

	//GetExplain 获取服务说明
	GetExplain() string

	//获取集群名称
	GetClusterName() string

	//获取分组名称
	GetGroupName() string

	//GetRPCPort 获取RPC端口
	GetRPCPort() int

	//GetHTTPPort 获取HTTP端口
	GetHTTPPort() int

	//GetDiscovery 获取服务发现配置
	GetDiscovery() []string

	//GetRedis 获取redis配置
	GetRedis() []string

	//GetEtcd 获取etcd配置
	GetEtcd() []string

	//GetKafka 获取kfaka地址
	GetKafka() []string

	//GetNats 获取nats地址
	GetNats() []string

	//GetRocketMq 获取rocketmq地址
	GetRocketMq() []string

	//GetNsq 获取nsq地址
	GetNsq() []string

	//GetReadMysql 获取ReadMysql地址
	GetReadMysql() *config.MysqlCfg

	//GetWriteMysql 获取GetWriteMysql地址
	GetWriteMysql() *config.MysqlCfg

	//GetHost 获取本机地址配置
	GetHost() string

	//GetRunmode 获取runmode地址配置
	GetRunmode() string

	//GetJaeger 获取链路追踪地址
	GetJaeger() string

	//GetModel 获取模式
	GetModel() string

	//GetMaxServiceLimitRequest 获取服务器最大的限流数
	GetMaxServiceLimitRequest() int

	//GetMaxClientLimitRequest 获取客户端最大的限流数
	GetMaxClientLimitRequest() int
}
