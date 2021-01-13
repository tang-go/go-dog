package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"

	"github.com/tang-go/go-dog/lib/net"
	"github.com/tang-go/go-dog/log"

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
	//Etcd地址
	Etcd []string `json:"etcd"`
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

//GetEtcd 获取etcd配置
func (c *Config) GetEtcd() []string {
	return c.Etcd
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

//GetJaeger 获取链路追踪地址
func (c *Config) GetJaeger() string {
	return c.Jaeger
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
	//Discovery 服务发现
	discovery := os.Getenv("DISCOVERY")
	if discovery != "" {
		var array []string
		if err := json.Unmarshal([]byte(discovery), &array); err != nil {
			c.Discovery = array
		}
	}
	//Redis地址
	redis := os.Getenv("REDIS")
	if redis != "" {
		var array []string
		if err := json.Unmarshal([]byte(redis), &array); err != nil {
			c.Redis = array
		}
	}
	//Etcd地址
	etcd := os.Getenv("ETCD")
	if etcd != "" {
		var array []string
		if err := json.Unmarshal([]byte(etcd), &array); err != nil {
			c.Etcd = array
		}
	}
	//Kafka地址
	kafka := os.Getenv("KAFKA")
	if kafka != "" {
		var array []string
		if err := json.Unmarshal([]byte(kafka), &array); err != nil {
			c.Kafka = array
		}
	}
	//Nats地址
	nats := os.Getenv("NATS")
	if nats != "" {
		var array []string
		if err := json.Unmarshal([]byte(nats), &array); err != nil {
			c.Nats = array
		}
	}
	//RocketMq地址
	rocketMq := os.Getenv("ROCKETMQ")
	if rocketMq != "" {
		var array []string
		if err := json.Unmarshal([]byte(rocketMq), &array); err != nil {
			c.RocketMq = array
		}
	}
	//nsq地址
	nsq := os.Getenv("NSQ")
	if nsq != "" {
		var array []string
		if err := json.Unmarshal([]byte(nsq), &array); err != nil {
			c.Nsq = array
		}
	}
	//Jaeger 链路追踪地址
	jaeger := os.Getenv("JAEGER")
	if jaeger != "" {
		c.Jaeger = jaeger
	}
	//读数据库
	{
		readMysqlIP := os.Getenv("READ_MYSQL_IP")
		if readMysqlIP != "" {
			c.ReadMysql.DbIP = readMysqlIP
		}
		readMysqlPwd := os.Getenv("READ_MYSQL_PWD")
		if readMysqlPwd != "" {
			c.ReadMysql.DbPWd = readMysqlPwd
		}
		readMysqlName := os.Getenv("READ_MYSQL_NAME")
		if readMysqlName != "" {
			c.ReadMysql.DbName = readMysqlName
		}
		readMysqlUser := os.Getenv("READ_MYSQL_USER")
		if readMysqlUser != "" {
			c.ReadMysql.DbUser = readMysqlUser
		}
		readMysqlMaxIdleConns := os.Getenv("READ_MYSQL_MAX_IDLE")
		if readMysqlMaxIdleConns != "" {
			maxIdleConns, err := strconv.Atoi(readMysqlMaxIdleConns)
			if err != nil {
				c.ReadMysql.MaxIdleConns = maxIdleConns
			}
		}
		readMysqlMaxOpenConns := os.Getenv("READ_MYSQL_MAX_OPEN")
		if readMysqlMaxOpenConns != "" {
			maxOpenConns, err := strconv.Atoi(readMysqlMaxOpenConns)
			if err != nil {
				c.ReadMysql.MaxOpenConns = maxOpenConns
			}
		}
		readMysqlOpenLog := os.Getenv("READ_MYSQL_OPEN_LOG")
		if readMysqlOpenLog != "" {
			openLog, err := strconv.ParseBool(readMysqlOpenLog)
			if err != nil {
				c.ReadMysql.OpenLog = openLog
			}
		}
	}
	//写数据库
	{
		writeMysqlIP := os.Getenv("WRITE_MYSQL_IP")
		if writeMysqlIP != "" {
			c.WriteMysql.DbIP = writeMysqlIP
		}
		writeMysqlPwd := os.Getenv("WRITE_MYSQL_PWD")
		if writeMysqlPwd != "" {
			c.WriteMysql.DbPWd = writeMysqlPwd
		}
		writeMysqlName := os.Getenv("WRITE_MYSQL_NAME")
		if writeMysqlName != "" {
			c.WriteMysql.DbName = writeMysqlName
		}
		writeMysqlUser := os.Getenv("WRITE_MYSQL_USER")
		if writeMysqlUser != "" {
			c.WriteMysql.DbUser = writeMysqlUser
		}
		writeMysqlMaxIdleConns := os.Getenv("WRITE_MYSQL_MAX_IDLE")
		if writeMysqlMaxIdleConns != "" {
			maxIdleConns, err := strconv.Atoi(writeMysqlMaxIdleConns)
			if err != nil {
				c.WriteMysql.MaxIdleConns = maxIdleConns
			}
		}
		writeMysqlMaxOpenConns := os.Getenv("WRITE_MYSQL_MAX_OPEN")
		if writeMysqlMaxOpenConns != "" {
			maxOpenConns, err := strconv.Atoi(writeMysqlMaxOpenConns)
			if err != nil {
				c.WriteMysql.MaxOpenConns = maxOpenConns
			}
		}
		writeMysqlOpenLog := os.Getenv("WRITE_MYSQL_OPEN_LOG")
		if writeMysqlOpenLog != "" {
			openLog, err := strconv.ParseBool(writeMysqlOpenLog)
			if err != nil {
				c.WriteMysql.OpenLog = openLog
			}
		}
	}

	//设置运行模式
	runmode := os.Getenv("RUNMODE")
	if runmode != "" {
		c.Runmode = runmode
	}
	//初始化日志
	switch c.GetRunmode() {
	case "panic":
		log.SetLevel(log.PanicLevel)
		break
	case "fatal":
		log.SetLevel(log.FatalLevel)
		break
	case "error":
		log.SetLevel(log.ErrorLevel)
		break
	case "warn":
		log.SetLevel(log.WarnLevel)
		break
	case "info":
		log.SetLevel(log.InfoLevel)
		break
	case "debug":
		log.SetLevel(log.DebugLevel)
		break
	case "trace":
		log.SetLevel(log.TraceLevel)
		break
	default:
		log.SetLevel(log.TraceLevel)
		break
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
	fmt.Println("### Etcd:         ", c.Etcd)
	fmt.Println("### Kafka:        ", c.Kafka)
	fmt.Println("### Nats:         ", c.Nats)
	fmt.Println("### RocketMq:     ", c.RocketMq)
	fmt.Println("### Nsq:          ", c.Nsq)
	fmt.Println("### ReadMysql:    ", c.ReadMysql)
	fmt.Println("### WriteMysql:   ", c.WriteMysql)
	fmt.Println("### Jaeger:       ", c.Jaeger)
	fmt.Println("### Host:         ", c.Host)
	fmt.Println("### RunMode:      ", c.Runmode)
	log.Traceln("日志初始化完成")
	return c
}
