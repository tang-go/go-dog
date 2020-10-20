package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/tang-go/go-dog/lib/net"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/sipt/GoJsoner"
)

var (
	configpath string
)

func init() {
	flag.StringVar(&configpath, "config", "./config/config.json", "config配置路径")
}

//Config 配置
type Config struct {
	//服务名称
	ServerName string `json:"server_name"`
	//服务说明
	Explain string `json:"explain"`
	//使用端口号
	Port int `json:"port"`
	//Discovery 服务发现
	Discovery []string `json:"discovery"`
	//Redis地址
	Redis []string `json:"redis"`
	//Kafka地址
	Kafka []string `json:"kafka"`
	//Nats地址
	Nats []string `json:"nats"`
	//RocketMq地址
	RocketMq []string `json:"rocket_mq"`
	//nsq地址
	Nsq []string `json:"nsq"`
	//Jaeger 链路追踪地址
	Jaeger string `json:"jaeger"`
	//读数据库地址
	ReadMysql *MysqlCfg `json:"read_mysql"`
	//写数据库地址
	WriteMysql *MysqlCfg `json:"write_mysql"`
	//本机地址
	Host string `json:"host"`
	//运行日志等级 panic fatal error warn info debug trace
	Runmode string `json:"runmode"`
}

//MysqlCfg mysql配置
type MysqlCfg struct {
	//数据库地址
	DbIP string `json:"db_ip"`
	//数据库密码
	DbPWd string `json:"db_pwd"`
	//数据库名称
	DbName string `json:"db_name"`
	//数据库用户
	DbUser string `json:"db_user"`
	//最大空闲连接数
	MaxIdleConns int `json:"max_idle_conns"`
	//最大连接数
	MaxOpenConns int `json:"max_open_conns"`
	//链接可重用时间
	ConnMaxLifetime int `json:"conn_max_lifetime"`
	//日志开关
	OpenLog bool `json:"open_log"`
}

//GetServerName 获取服务名称
func (c *Config) GetServerName() string {
	return c.ServerName
}

//GetPort 获取端口
func (c *Config) GetPort() int {
	return c.Port
}

//GetExplain 获取服务说明
func (c *Config) GetExplain() string {
	return c.Explain
}

//GetDiscovery 获取服务发现配置
func (c *Config) GetDiscovery() []string {
	return c.Discovery
}

//GetRedis 获取redis配置
func (c *Config) GetRedis() []string {
	return c.Redis
}

//GetKafka 获取kfaka地址
func (c *Config) GetKafka() []string {
	return c.Kafka
}

//GetNats 获取nats地址
func (c *Config) GetNats() []string {
	return c.Nats
}

//GetRocketMq 获取rocketmq地址
func (c *Config) GetRocketMq() []string {
	return c.RocketMq
}

//GetNsq 获取nsq地址
func (c *Config) GetNsq() []string {
	return c.Nsq
}

//GetReadMysql 获取ReadMysql地址
func (c *Config) GetReadMysql() *MysqlCfg {
	return c.ReadMysql
}

//GetWriteMysql 获取GetWriteMysql地址
func (c *Config) GetWriteMysql() *MysqlCfg {
	return c.WriteMysql
}

//GetHost 获取本机地址配置
func (c *Config) GetHost() string {
	return c.Host
}

//GetRunmode 获取runmode地址配置
func (c *Config) GetRunmode() string {
	return c.Runmode
}

//NewConfig 初始化Config
func NewConfig() *Config {
	c := new(Config)
	//从文件读取json文件并且解析
	flag.Parse()
	s := os.Getenv("config")
	if s == "" {
		gameConfigData, err := ioutil.ReadFile(configpath)
		if err != nil {
			panic(err.Error())
		}
		gameConfigResult, err := GoJsoner.Discard(string(gameConfigData))
		if err != nil {
			panic(err.Error())
		}
		err = json.Unmarshal([]byte(gameConfigResult), c)
		if err != nil {
			panic(err.Error())
		}
	} else {
		gameConfigResult, err := GoJsoner.Discard(s)
		if err != nil {
			panic(err.Error())
		}
		err = json.Unmarshal([]byte(gameConfigResult), c)
		if err != nil {
			panic(err.Error())
		}
	}

	host := os.Getenv("HOST")
	if host != "" {
		c.Host = host
	} else {
		if c.Host == "" {
			host, err := net.ExternalIP()
			if err != nil {
				panic(err.Error())
			}
			c.Host = host.String()
		}
	}
	//先看环境变量是否有端口号
	port := os.Getenv("PORT")
	if port != "" {
		p, err := strconv.Atoi(port)
		if err != nil {
			c.Port = p
		}
	}
	if c.Port <= 0 {
		//获取随机端口
		p, err := net.GetFreePort()
		if err != nil {
			panic(err.Error())
		}
		c.Port = p
	}
	//设置运行模式
	runmode := os.Getenv("RUNMODE")
	if runmode != "" {
		c.Runmode = runmode
	}
	fmt.Println("************************************************")
	fmt.Println("*                                              *")
	fmt.Println("*             	   Cfg  Init                    *")
	fmt.Println("*                                              *")
	fmt.Println("************************************************")
	fmt.Println("### ServerName:   ", c.ServerName)
	fmt.Println("### Port:         ", c.Port)
	fmt.Println("### Discovery:    ", c.Discovery)
	fmt.Println("### Redis:        ", c.Redis)
	fmt.Println("### Kafka:        ", c.Kafka)
	fmt.Println("### Nats:         ", c.Nats)
	fmt.Println("### RocketMq:     ", c.RocketMq)
	fmt.Println("### Nsq:          ", c.Nsq)
	fmt.Println("### ReadMysql:    ", c.ReadMysql)
	fmt.Println("### WriteMysql:   ", c.WriteMysql)
	fmt.Println("### Jaeger:       ", c.Jaeger)
	fmt.Println("### Host:         ", c.Host)
	fmt.Println("### RunMode:      ", c.Runmode)
	return c
}
