package plugins

import "github.com/tang-go/go-dog/pkg/config"

//Cfg 配置接口
type Cfg interface {
	//GetServerName 获取服务名称
	GetServerName() string

	//GetExplain 获取服务说明
	GetExplain() string

	//GetPort 获取端口
	GetPort() int

	//GetDiscovery 获取服务发现配置
	GetDiscovery() []string

	//GetRedis 获取redis配置
	GetRedis() []string

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
}
