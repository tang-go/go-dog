package client

import (
	customerror "go-dog/error"
	"go-dog/internal/codec"
	"go-dog/internal/config"
	"go-dog/internal/discovery"
	"go-dog/internal/fusing"
	"go-dog/internal/limit"
	"go-dog/internal/selector"
	"go-dog/pkg/log"
	"go-dog/pkg/recover"
	"go-dog/plugins"
	"go-dog/serviceinfo"
	"sync"
	"time"
)

const (
	_MaxClientRequestCount = 100000
)

//Client 客户端
type Client struct {
	cfg           plugins.Cfg
	codec         plugins.Codec
	discovery     plugins.Discovery
	fusing        plugins.Fusing
	selector      plugins.Selector
	limit         plugins.Limit
	managerclient *ManagerClient
	wait          sync.WaitGroup
}

//NewClient  新建一个客户端
func NewClient(discoveryTTL int64, param ...interface{}) plugins.Client {
	client := new(Client)
	for _, plugin := range param {
		if cfg, ok := plugin.(plugins.Cfg); ok {
			client.cfg = cfg
		}
		if discovery, ok := plugin.(plugins.Discovery); ok {
			client.discovery = discovery
		}
		if fusing, ok := plugin.(plugins.Fusing); ok {
			client.fusing = fusing
		}
		if selector, ok := plugin.(plugins.Selector); ok {
			client.selector = selector
		}
		if limit, ok := plugin.(plugins.Limit); ok {
			client.limit = limit
		}
		if codec, ok := plugin.(plugins.Codec); ok {
			client.codec = codec
		}
	}
	if client.cfg == nil {
		//默认配置
		client.cfg = config.NewConfig()
	}
	if client.discovery == nil {
		//使用默认服务发现中心
		client.discovery = discovery.NewEtcdDiscovery(client.cfg.GetEtcd(), discoveryTTL)
	}
	if client.fusing == nil {
		//使用默认的熔断插件
		client.fusing = fusing.NewFusing(time.Duration(discoveryTTL) * time.Second)
	}
	if client.selector == nil {
		//使用默认的选择器
		client.selector = selector.NewSelector()
	}
	if client.limit == nil {
		//使用默认的流量限制
		client.limit = limit.NewLimit(_MaxClientRequestCount)
	}
	if client.codec == nil {
		//使用默认的参数编码
		client.codec = codec.NewCodec()
	}
	//初始化日志
	switch client.cfg.GetRunmode() {
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
	client.managerclient = NewManagerClient(client.codec)
	//开始RPC监听服务上下线
	client.discovery.RegRPCServiceOnlineNotice(client.ServiceOnlineNotice)
	client.discovery.RegRPCServiceOfflineNotice(client.ServiceOfflineNotice)
	client.discovery.WatchRPCService()
	return client
}

//GetCfg 获取配置
func (c *Client) GetCfg() plugins.Cfg {
	return c.cfg
}

//GetDiscovery 获取服务发现
func (c *Client) GetDiscovery() plugins.Discovery {
	return c.discovery
}

//GetFusing 获取熔断插件
func (c *Client) GetFusing() plugins.Fusing {
	return c.fusing
}

//GetLimit 获取限流插件
func (c *Client) GetLimit() plugins.Limit {
	return c.limit
}

//GetCodec 获取编码插件
func (c *Client) GetCodec() plugins.Codec {
	return c.codec
}

//GetAllService 获取所有服务
func (c *Client) GetAllService() (services []*serviceinfo.ServiceInfo) {
	return c.selector.GetAllService()
}

//ServiceOnlineNotice 服务上线
func (c *Client) ServiceOnlineNotice(key string, info *serviceinfo.ServiceInfo) {
	c.selector.AddService(key, info)
}

//ServiceOfflineNotice 服务下线
func (c *Client) ServiceOfflineNotice(key string) {
	c.selector.DelService(key)
	c.managerclient.DelClient(key)
}

//Call 调用函数
func (c *Client) Call(ctx plugins.Context, mode plugins.Mode, name string, method string, args interface{}, reply interface{}) error {
	defer recover.Recover()
	if c.limit.IsLimit() {
		return customerror.EnCodeError(customerror.ClientLimitError, "超过了每秒最大流量")
	}
	c.wait.Add(1)
	defer c.wait.Done()
	switch mode {
	//随机模式
	case plugins.RandomMode:
		service, err := c.selector.RandomMode(c.fusing, name, method)
		if err != nil {
			return err
		}
		client, err := c.managerclient.GetClient(service)
		if err == nil {
			//请求统计添加
			c.fusing.AddMethod(service.Key, method)
			//客户端发起请求
			err := client.Call(ctx, name, method, args, reply)
			if err != nil {
				//添加错误
				c.fusing.AddErrorMethod(service.Key, method, err)
				return err
			}
			return nil
		}
		return customerror.EnCodeError(customerror.InternalServerError, "没有服务可用")
	//遍历模式
	case plugins.RangeMode:
		var e error = customerror.EnCodeError(customerror.InternalServerError, "没有服务可用")
		c.selector.RangeMode(c.fusing, name, method, func(service *serviceinfo.ServiceInfo) bool {
			client, err := c.managerclient.GetClient(service)
			if err != nil {
				e = err
				return false
			}
			//请求统计添加
			c.fusing.AddMethod(service.Key, method)
			err = client.Call(ctx, name, method, args, reply)
			if err != nil {
				//添加错误
				c.fusing.AddErrorMethod(service.Key, method, err)
				e = err
				return false
			}
			return true
		})
		return e
	//hash模式
	case plugins.HashMode:
		service, err := c.selector.HashMode(c.fusing, name, method)
		if err != nil {
			return err
		}
		client, err := c.managerclient.GetClient(service)
		if err == nil {
			//请求统计添加
			c.fusing.AddMethod(service.Key, method)
			//客户端发起请求
			err := client.Call(ctx, name, method, args, reply)
			if err != nil {
				//添加错误
				c.fusing.AddErrorMethod(service.Key, method, err)
				return err
			}
			return nil
		}
		return customerror.EnCodeError(customerror.InternalServerError, "没有服务可用")
	//默认方式
	default:
		service, err := c.selector.Custom(c.fusing, name, method)
		if err != nil {
			return err
		}
		client, err := c.managerclient.GetClient(service)
		if err == nil {
			//请求统计添加
			c.fusing.AddMethod(service.Key, method)
			//客户端发起请求
			err := client.Call(ctx, name, method, args, reply)
			if err != nil {
				//添加错误
				c.fusing.AddErrorMethod(service.Key, method, err)
				return err
			}
			return nil
		}
		return customerror.EnCodeError(customerror.InternalServerError, "没有服务可用")
	}
}

//SendRequest 发生请求
func (c *Client) SendRequest(ctx plugins.Context, mode plugins.Mode, name string, method string, args []byte) (reply []byte, e error) {
	defer recover.Recover()
	if c.limit.IsLimit() {
		return nil, customerror.EnCodeError(customerror.ClientLimitError, "超过了每秒最大流量")
	}
	c.wait.Add(1)
	defer c.wait.Done()
	switch mode {
	//随机模式
	case plugins.RandomMode:
		service, err := c.selector.RandomMode(c.fusing, name, method)
		if err != nil {
			return nil, err
		}
		client, err := c.managerclient.GetClient(service)
		if err == nil {
			//请求统计添加
			c.fusing.AddMethod(service.Key, method)
			//客户端发起请求
			res, err := client.SendRequest(ctx, name, method, args)
			if err != nil {
				//添加错误
				c.fusing.AddErrorMethod(service.Key, method, err)
				return nil, err
			}
			return res, nil
		}
		return nil, customerror.EnCodeError(customerror.InternalServerError, "没有服务可用")
	//遍历模式
	case plugins.RangeMode:
		var e error = customerror.EnCodeError(customerror.InternalServerError, "没有服务可用")
		var res []byte
		c.selector.RangeMode(c.fusing, name, method, func(service *serviceinfo.ServiceInfo) bool {
			client, err := c.managerclient.GetClient(service)
			if err != nil {
				e = err
				return false
			}
			//请求统计添加
			c.fusing.AddMethod(service.Key, method)
			res, err = client.SendRequest(ctx, name, method, args)
			if err != nil {
				//添加错误
				c.fusing.AddErrorMethod(service.Key, method, err)
				e = err
				return false
			}
			return true
		})
		return res, e
	//hash模式
	case plugins.HashMode:
		service, err := c.selector.HashMode(c.fusing, name, method)
		if err != nil {
			return nil, err
		}
		client, err := c.managerclient.GetClient(service)
		if err == nil {
			//请求统计添加
			c.fusing.AddMethod(service.Key, method)
			//客户端发起请求
			res, err := client.SendRequest(ctx, name, method, args)
			if err != nil {
				//添加错误
				c.fusing.AddErrorMethod(service.Key, method, err)
				return nil, err
			}
			return res, nil
		}
		return nil, customerror.EnCodeError(customerror.InternalServerError, "没有服务可用")
	//默认方式
	default:
		service, err := c.selector.Custom(c.fusing, name, method)
		if err != nil {
			return nil, err
		}
		client, err := c.managerclient.GetClient(service)
		if err == nil {
			//请求统计添加
			c.fusing.AddMethod(service.Key, method)
			//客户端发起请求
			res, err := client.SendRequest(ctx, name, method, args)
			if err != nil {
				//添加错误
				c.fusing.AddErrorMethod(service.Key, method, err)
				return nil, err
			}
			return res, nil
		}
		return nil, customerror.EnCodeError(customerror.InternalServerError, "没有服务可用")
	}
}

//Close 关闭
func (c *Client) Close() {
	c.managerclient.Close()
	c.wait.Wait()
	c.discovery.Close()
	c.fusing.Close()
	c.limit.Close()
}
