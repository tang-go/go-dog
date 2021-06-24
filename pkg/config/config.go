package config

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/sipt/GoJsoner"
	"github.com/tang-go/go-dog/lib/net"
	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/nacos"
)

const (
	localModel = "local"
	nacosModel = "nacos"
)
const (
	ConsulDiscoveryModel = "consul"
	NacosDiscoveryModel  = "nacos"
)

var (
	configpath     string
	modle          string
	discoveryModel string
)

const (
	_MaxClientRequestCount  int = 100000
	_MaxServiceRequestCount int = 10000
)

//NacosConfig 配置
type NacosConfig struct {
	//命名空间 空为默认
	Namespace string `json:"namespace"`
	//用户名称
	Username string `json:"username"`
	//密码
	Password string `json:"passwork"`
	//关心配置的DataID
	DataID string `json:"data_id"`
	//关心配置的组
	Group string `json:"group"`
	//Nacos 地址
	Address []nacos.Address `json:"address"`
}

func init() {
	flag.StringVar(&configpath, "c", "./config/config.json", "config配置路径")
	flag.StringVar(&discoveryModel, "d", "consul", "nacos nacos服务发型模式;consul consul服务发型模式")
	flag.StringVar(&modle, "m", "loacl", "loacl 本地配置模式;nacos nacos配置模式")

}

//Config 配置
type Config struct {
	//服务名称
	ServerName string `json:"server_name"`
	//服务说明
	Explain string `json:"explain"`
	//ClusterName 集群名称
	ClusterName string `json:"cluster_name"`
	//GroupName 分组名称
	GroupName string `json:"group_name"`
	//RPC使用的端口
	RPCPort int `json:"rpc_port"`
	//HTTP使用的端口
	HTTPPort int `json:"http_port"`
	//Consul Consul地址
	Consul string `json:"consul"`
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
	//Nacos 地址
	Nacos *NacosConfig `json:"nacos"`
	//本机地址
	Host string `json:"host"`
	//运行日志等级 panic fatal error warn info debug trace
	Runmode string `json:"runmode"`
	//服务端最大接受的处理数量
	MaxServiceLimitRequest int `json:"max_service_limit_request"`
	//客户端最大的请求数量
	MaxClientLimitRequest int `json:"max_client_limit_request"`
	//模式
	Model string `json:"-"`
	//服务发型模式
	DiscoveryModel string `json:"-"`
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

//GetClusterName 获取集群名称
func (c *Config) GetClusterName() string {
	return c.ClusterName
}

//GetGroupName 获取分组名称
func (c *Config) GetGroupName() string {
	return c.GroupName
}

//GetServerName 获取服务名称
func (c *Config) GetServerName() string {
	return c.ServerName
}

//GetRPCPort 获取RPC端口
func (c *Config) GetRPCPort() int {
	return c.RPCPort
}

//GetHTTPPort 获取HTTP端口
func (c *Config) GetHTTPPort() int {
	return c.HTTPPort
}

//GetExplain 获取服务说明
func (c *Config) GetExplain() string {
	return c.Explain
}

//GetConsul 获取consul
func (c *Config) GetConsul() string {
	return c.Consul
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

//GetModel 获取模式
func (c *Config) GetModel() string {
	return c.Model
}

//GetDiscoveryModel 获取服务发现模型
func (c *Config) GetDiscoveryModel() string {
	return c.DiscoveryModel
}

//GetMaxServiceLimitRequest 获取服务器最大的限流数
func (c *Config) GetMaxServiceLimitRequest() int {
	return c.MaxServiceLimitRequest
}

//GetMaxClientLimitRequest 获取客户端最大的限流数
func (c *Config) GetMaxClientLimitRequest() int {
	return c.MaxClientLimitRequest
}

//NewConfig 初始化Config
func NewConfig() *Config {
	//从文件读取json文件并且解析
	flag.Parse()
	c := new(Config)
	c.Model = modle
	c.DiscoveryModel = discoveryModel
	c.initCfgModel()
	c.initEnv()
	fmt.Println("************************************************")
	fmt.Println("*                                              *")
	fmt.Println("*             	   Cfg  Init                    *")
	fmt.Println("*                                              *")
	fmt.Println("************************************************")
	fmt.Println("### Model:        ", c.Model)
	fmt.Println("### ServerName:   ", c.ServerName)
	fmt.Println("### ClusterName:  ", c.ClusterName)
	fmt.Println("### GroupName:    ", c.GroupName)
	fmt.Println("### Explain:      ", c.Explain)
	fmt.Println("### RPCPort:      ", c.RPCPort)
	fmt.Println("### HTTPPort:     ", c.HTTPPort)
	fmt.Println("### Consul:       ", c.Consul)
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
	fmt.Println("### ServiceLimit: ", c.MaxServiceLimitRequest)
	fmt.Println("### ClientLimit:  ", c.MaxClientLimitRequest)
	fmt.Println("### RunMode:      ", c.Runmode)
	log.Traceln("日志初始化完成")
	return c
}

func (c *Config) initLog() {
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
}

func (c *Config) initCfgModel() {
	switch c.Model {
	case nacosModel:
		//远程nacos配置模型
		nacosConfig := new(NacosConfig)
		configData, err := ioutil.ReadFile(configpath)
		if err != nil {
			panic(err.Error())
		}
		configResult, err := GoJsoner.Discard(string(configData))
		if err != nil {
			panic(err.Error())
		}
		err = json.Unmarshal([]byte(configResult), nacosConfig)
		if err != nil {
			panic(err.Error())
		}
		//从环境变量获取
		s := os.Getenv("config")
		if s != "" {
			configResult, err := GoJsoner.Discard(s)
			if err != nil {
				panic(err.Error())
			}
			err = json.Unmarshal([]byte(configResult), nacosConfig)
			if err != nil {
				panic(err.Error())
			}
		}
		if namespace := os.Getenv("NAMESPACE"); namespace != "" {
			nacosConfig.Namespace = namespace
		}
		if username := os.Getenv("USERNAME"); username != "" {
			nacosConfig.Username = username
		}
		if password := os.Getenv("PASSWORK"); password != "" {
			nacosConfig.Password = password
		}
		if dataID := os.Getenv("DATA_ID"); dataID != "" {
			nacosConfig.DataID = dataID
		}
		if group := os.Getenv("GROUP"); group != "" {
			nacosConfig.Group = group
		}
		//初始化nacos
		nacos.Init(nacosConfig.Namespace, nacosConfig.Username, nacosConfig.Password, nacosConfig.Address)
		//初始化配置
		cfg, err := nacos.GetConfig().GetConfig(nacosConfig.DataID, nacosConfig.Group)
		if err != nil {
			panic(err)
		}
		//解析配置
		gameConfigResult, err := GoJsoner.Discard(cfg)
		if err != nil {
			panic(err.Error())
		}
		err = json.Unmarshal([]byte(gameConfigResult), c)
		if err != nil {
			panic(err.Error())
		}
	default:
		//默认本地模式
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
		s := os.Getenv("config")
		if s != "" {
			gameConfigResult, err := GoJsoner.Discard(s)
			if err != nil {
				panic(err.Error())
			}
			err = json.Unmarshal([]byte(gameConfigResult), c)
			if err != nil {
				panic(err.Error())
			}
		}
	}
}

func (c *Config) initEnv() {
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
	//服务器限流
	maxServiceLimitRequest := os.Getenv("MAX_SERVICE_LIMIT_REQUEST")
	if maxServiceLimitRequest != "" {
		p, err := strconv.Atoi(maxServiceLimitRequest)
		if err != nil {
			panic(err.Error())
		}
		c.MaxServiceLimitRequest = p
	}
	if c.MaxServiceLimitRequest <= 0 {
		c.MaxServiceLimitRequest = _MaxServiceRequestCount
	}
	//客户端限流
	maxClientLimitRequest := os.Getenv("MAX_CLIENT_LIMIT_REQUEST")
	if maxClientLimitRequest != "" {
		p, err := strconv.Atoi(maxClientLimitRequest)
		if err != nil {
			panic(err.Error())
		}
		c.MaxClientLimitRequest = p
	}
	if c.MaxClientLimitRequest <= 0 {
		c.MaxClientLimitRequest = _MaxClientRequestCount
	}
	//先看环境变量是否有端口号
	rpcport := os.Getenv("RPC_PORT")
	if rpcport != "" {
		p, err := strconv.Atoi(rpcport)
		if err != nil {
			panic(err.Error())
		}
		c.RPCPort = p
	}
	if c.RPCPort <= 0 {
		//获取随机端口
		p, err := net.GetFreePort()
		if err != nil {
			panic(err.Error())
		}
		c.RPCPort = p
	}
	//先看环境变量是否有端口号
	httpport := os.Getenv("HTTP_PORT")
	if httpport != "" {
		p, err := strconv.Atoi(httpport)
		if err != nil {
			panic(err.Error())
		}
		c.HTTPPort = p
	}
	if c.HTTPPort <= 0 {
		//获取随机端口
		p, err := net.GetFreePort()
		if err != nil {
			panic(err.Error())
		}
		c.HTTPPort = p
	}
	//consul consul地址
	consul := os.Getenv("CONSUL")
	if consul != "" {
		c.Consul = consul
	}
	//Redis地址
	redis := os.Getenv("REDIS")
	if redis != "" {
		array := strings.Split(redis, ",")
		c.Redis = array
	}
	//Etcd地址
	etcd := os.Getenv("ETCD")
	if etcd != "" {
		array := strings.Split(etcd, ",")
		c.Etcd = array
	}
	//Kafka地址
	kafka := os.Getenv("KAFKA")
	if kafka != "" {
		array := strings.Split(kafka, ",")
		c.Kafka = array
	}
	//Nats地址
	nats := os.Getenv("NATS")
	if nats != "" {
		array := strings.Split(nats, ",")
		c.Nats = array
	}
	//RocketMq地址
	rocketMq := os.Getenv("ROCKETMQ")
	if rocketMq != "" {
		array := strings.Split(rocketMq, ",")
		c.RocketMq = array
	}
	//nsq地址
	nsq := os.Getenv("NSQ")
	if nsq != "" {
		array := strings.Split(nsq, ",")
		c.Nsq = array
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
				panic(err)
			}
			c.ReadMysql.MaxIdleConns = maxIdleConns
		}
		readMysqlMaxOpenConns := os.Getenv("READ_MYSQL_MAX_OPEN")
		if readMysqlMaxOpenConns != "" {
			maxOpenConns, err := strconv.Atoi(readMysqlMaxOpenConns)
			if err != nil {
				panic(err)
			}
			c.ReadMysql.MaxOpenConns = maxOpenConns
		}
		readMysqlOpenLog := os.Getenv("READ_MYSQL_OPEN_LOG")
		if readMysqlOpenLog != "" {
			openLog, err := strconv.ParseBool(readMysqlOpenLog)
			if err != nil {
				panic(err)
			}
			c.ReadMysql.OpenLog = openLog
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
				panic(err)
			}
			c.WriteMysql.MaxIdleConns = maxIdleConns
		}
		writeMysqlMaxOpenConns := os.Getenv("WRITE_MYSQL_MAX_OPEN")
		if writeMysqlMaxOpenConns != "" {
			maxOpenConns, err := strconv.Atoi(writeMysqlMaxOpenConns)
			if err != nil {
				panic(err)
			}
			c.WriteMysql.MaxOpenConns = maxOpenConns
		}
		writeMysqlOpenLog := os.Getenv("WRITE_MYSQL_OPEN_LOG")
		if writeMysqlOpenLog != "" {
			openLog, err := strconv.ParseBool(writeMysqlOpenLog)
			if err != nil {
				panic(err)
			}
			c.WriteMysql.OpenLog = openLog
		}
	}

	//设置集群名称
	clusterName := os.Getenv("CLUSTER_NAME")
	if clusterName != "" {
		c.ClusterName = clusterName
	}

	//设置服务分组
	groupName := os.Getenv("GROUP_NAME")
	if groupName != "" {
		c.GroupName = groupName
	}

	//设置运行模式
	runmode := os.Getenv("RUNMODE")
	if runmode != "" {
		c.Runmode = runmode
	}
}
