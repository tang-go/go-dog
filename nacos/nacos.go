package nacos

import (
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

var gNacos *Nacos

//Nacos nacos配置中心
type Nacos struct {
	register  *Register
	discovery *Discovery
	config    *Config
}

//Address 地址
type Address struct {
	IP   string `json:"ip"`
	Port uint64 `json:"port"`
}

//初始化nacos
func Init(namespace, username, password string, address []Address) {
	gNacos = newNacos(namespace, username, password, address)
}

//newNacos 初始化私有nacos对象
func newNacos(namespace, username, password string, address []Address) *Nacos {
	sc := make([]constant.ServerConfig, 0)
	for _, add := range address {
		sc = append(sc, constant.ServerConfig{
			IpAddr:      add.IP,
			Port:        add.Port,
			ContextPath: "/nacos",
			Scheme:      "http",
		})
	}
	cc := &constant.ClientConfig{
		NamespaceId:          namespace,
		BeatInterval:         1 * 1000,
		NotLoadCacheAtStart:  true,
		LogDir:               "./log",
		CacheDir:             "./cache",
		RotateTime:           "1h",
		MaxAge:               3,
		LogLevel:             "error",
		Username:             username,
		Password:             password,
		UpdateCacheWhenEmpty: true,
	}
	//初始化一个服务发现的客户端
	inamingClient, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		panic(err.Error())
	}
	//初始化配置中心
	configClient, err := clients.NewConfigClient(
		vo.NacosClientParam{
			ClientConfig:  cc,
			ServerConfigs: sc,
		},
	)
	if err != nil {
		panic(err)
	}
	n := new(Nacos)
	n.discovery = newDiscovery(inamingClient)
	n.register = newRegister(inamingClient)
	n.config = newConfig(configClient)
	return n
}

//GetRegister 获取注册中心
func GetRegister() *Register {
	return gNacos.register
}

//GetDiscovery 获取服务发现
func GetDiscovery() *Discovery {
	return gNacos.discovery
}

//GetConfig 获取配置中心
func GetConfig() *Config {
	return gNacos.config
}
