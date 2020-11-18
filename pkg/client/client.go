package client

import (
	"fmt"
	"sync"
	"time"

	customerror "github.com/tang-go/go-dog/error"
	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/pkg/codec"
	"github.com/tang-go/go-dog/pkg/config"
	"github.com/tang-go/go-dog/pkg/context"
	"github.com/tang-go/go-dog/pkg/discovery"
	"github.com/tang-go/go-dog/pkg/fusing"
	"github.com/tang-go/go-dog/pkg/limit"
	"github.com/tang-go/go-dog/pkg/selector"
	"github.com/tang-go/go-dog/plugins"
	"github.com/tang-go/go-dog/recover"
	"github.com/tang-go/go-dog/serviceinfo"
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
func NewClient(param ...interface{}) plugins.Client {
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
		client.discovery = discovery.NewGoDogDiscovery(client.cfg.GetDiscovery())
	}
	if client.fusing == nil {
		//使用默认的熔断插件
		client.fusing = fusing.NewFusing(2 * time.Second)
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

//GetAllRPCService 获取所有RPC服务
func (c *Client) GetAllRPCService() (services []*serviceinfo.RPCServiceInfo) {
	return c.discovery.GetAllRPCService()
}

//GetAllAPIService 获取所有API服务
func (c *Client) GetAllAPIService() (services []*serviceinfo.APIServiceInfo) {
	return c.discovery.GetAllAPIService()
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
		service, err := c.selector.RandomMode(c.discovery, c.fusing, name, method)
		if err != nil {
			log.Errorln(err.Error())
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
				log.Errorln(err.Error())
				c.fusing.AddErrorMethod(service.Key, method, err)
				return err
			}
			if err != nil {
				log.Errorln(err.Error())
			}
			return nil
		}
		return customerror.EnCodeError(customerror.InternalServerError, "没有服务可用")
	//遍历模式
	case plugins.RangeMode:
		var e error = customerror.EnCodeError(customerror.InternalServerError, "没有服务可用")
		e = c.selector.RangeMode(c.discovery, c.fusing, name, method, func(service *serviceinfo.RPCServiceInfo) bool {
			client, err := c.managerclient.GetClient(service)
			if err != nil {
				e = err
				log.Errorln(err.Error())
				return false
			}
			//请求统计添加
			c.fusing.AddMethod(service.Key, method)
			err = client.Call(ctx, name, method, args, reply)
			if err != nil {
				//添加错误
				log.Errorln(err.Error())
				c.fusing.AddErrorMethod(service.Key, method, err)
				e = err
				return false
			}
			return true
		})
		if e != nil {
			log.Errorln(e.Error())
		}
		return e
	//hash模式
	case plugins.HashMode:
		service, err := c.selector.HashMode(c.discovery, c.fusing, name, method)
		if err != nil {
			log.Errorln(err.Error())
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
				log.Errorln(err.Error())
				c.fusing.AddErrorMethod(service.Key, method, err)
				return err
			}
			return nil
		}
		if err != nil {
			log.Errorln(err.Error())
		}
		return customerror.EnCodeError(customerror.InternalServerError, "没有服务可用")
	//默认方式
	default:
		service, err := c.selector.Custom(c.discovery, c.fusing, name, method)
		if err != nil {
			log.Errorln(err.Error())
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
				log.Errorln(err.Error())
				c.fusing.AddErrorMethod(service.Key, method, err)
				return err
			}
			return nil
		}
		if err != nil {
			log.Errorln(err.Error())
		}
		return customerror.EnCodeError(customerror.InternalServerError, "没有服务可用")
	}
}

//SendRequest 发生请求
func (c *Client) SendRequest(ctx plugins.Context, mode plugins.Mode, name string, method string, code string, args []byte) (reply []byte, e error) {
	defer recover.Recover()
	if c.limit.IsLimit() {
		return nil, customerror.EnCodeError(customerror.ClientLimitError, "超过了每秒最大流量")
	}
	c.wait.Add(1)
	defer c.wait.Done()
	switch mode {
	//随机模式
	case plugins.RandomMode:
		service, err := c.selector.RandomMode(c.discovery, c.fusing, name, method)
		if err != nil {
			log.Errorln(err.Error())
			return nil, err
		}
		client, err := c.managerclient.GetClient(service)
		if err == nil {
			//请求统计添加
			c.fusing.AddMethod(service.Key, method)
			//客户端发起请求
			res, err := client.SendRequest(ctx, name, method, code, args)
			if err != nil {
				//添加错误
				log.Errorln(err.Error())
				c.fusing.AddErrorMethod(service.Key, method, err)
				return nil, err
			}
			return res, nil
		}
		if err != nil {
			log.Errorln(err.Error())
		}
		return nil, customerror.EnCodeError(customerror.InternalServerError, "没有服务可用")
	//遍历模式
	case plugins.RangeMode:
		var e error = customerror.EnCodeError(customerror.InternalServerError, "没有服务可用")
		var res []byte
		e = c.selector.RangeMode(c.discovery, c.fusing, name, method, func(service *serviceinfo.RPCServiceInfo) bool {
			client, err := c.managerclient.GetClient(service)
			if err != nil {
				e = err
				log.Errorln(err.Error())
				return false
			}
			//请求统计添加
			c.fusing.AddMethod(service.Key, method)
			res, err = client.SendRequest(ctx, name, method, code, args)
			if err != nil {
				//添加错误
				log.Errorln(err.Error())
				c.fusing.AddErrorMethod(service.Key, method, err)
				e = err
				return false
			}
			return true
		})
		if e != nil {
			log.Errorln(e.Error())
		}
		return res, e
	//hash模式
	case plugins.HashMode:
		service, err := c.selector.HashMode(c.discovery, c.fusing, name, method)
		if err != nil {
			return nil, err
		}
		client, err := c.managerclient.GetClient(service)
		if err == nil {
			//请求统计添加
			c.fusing.AddMethod(service.Key, method)
			//客户端发起请求
			res, err := client.SendRequest(ctx, name, method, code, args)
			if err != nil {
				//添加错误
				log.Errorln(err.Error())
				c.fusing.AddErrorMethod(service.Key, method, err)
				return nil, err
			}
			return res, nil
		}
		if err != nil {
			log.Errorln(err.Error())
		}
		return nil, customerror.EnCodeError(customerror.InternalServerError, "没有服务可用")
	//默认方式
	default:
		service, err := c.selector.Custom(c.discovery, c.fusing, name, method)
		if err != nil {
			log.Errorln(err.Error())
			return nil, err
		}
		client, err := c.managerclient.GetClient(service)
		if err == nil {
			//请求统计添加
			c.fusing.AddMethod(service.Key, method)
			//客户端发起请求
			res, err := client.SendRequest(ctx, name, method, code, args)
			if err != nil {
				//添加错误
				log.Errorln(err.Error())
				c.fusing.AddErrorMethod(service.Key, method, err)
				return nil, err
			}
			return res, nil
		}
		if err != nil {
			log.Errorln(err.Error())
		}
		return nil, customerror.EnCodeError(customerror.InternalServerError, "没有服务可用")
	}
}

//Broadcast 广播
func (c *Client) Broadcast(ctx plugins.Context, name string, method string, args interface{}, reply interface{}) error {
	defer recover.Recover()
	if c.limit.IsLimit() {
		return customerror.EnCodeError(customerror.ClientLimitError, "超过了每秒最大流量")
	}
	c.wait.Add(1)
	defer c.wait.Done()
	var e error = customerror.EnCodeError(customerror.InternalServerError, "没有服务可用")
	e = c.selector.RangeMode(c.discovery, c.fusing, name, method, func(service *serviceinfo.RPCServiceInfo) bool {
		client, err := c.managerclient.GetClient(service)
		if err != nil {
			e = err
			log.Errorln(err.Error())
			return false
		}
		//请求统计添加
		c.fusing.AddMethod(service.Key, method)
		err = client.Call(context.WithTimeout(ctx, int64(time.Second*5)), name, method, args, reply)
		if err != nil {
			//添加错误
			log.Errorln(err.Error())
			c.fusing.AddErrorMethod(service.Key, method, err)
			e = err
			return false
		}
		return false
	})
	if e != nil {
		log.Errorln(e.Error())
	}
	return e
}

//CallByAddress 指定地址调用
func (c *Client) CallByAddress(ctx plugins.Context, address string, name string, method string, args interface{}, reply interface{}) error {
	defer recover.Recover()
	if c.limit.IsLimit() {
		return customerror.EnCodeError(customerror.ClientLimitError, "超过了每秒最大流量")
	}
	c.wait.Add(1)
	defer c.wait.Done()

	service, err := c.selector.GetByAddress(c.discovery, address, c.fusing, name, method)
	if err != nil {
		log.Errorln(err.Error())
		return err
	}
	ctx.SetSource(fmt.Sprintf("%s:%d", c.cfg.GetHost(), c.cfg.GetPort()))
	client, err := c.managerclient.GetClient(service)
	if err == nil {
		//请求统计添加
		c.fusing.AddMethod(service.Key, method)
		//客户端发起请求
		err := client.Call(ctx, name, method, args, reply)
		if err != nil {
			//添加错误
			log.Errorln(err.Error())
			c.fusing.AddErrorMethod(service.Key, method, err)
			return err
		}
		return nil
	}
	if err != nil {
		log.Errorln(err.Error())
	}
	return customerror.EnCodeError(customerror.InternalServerError, "没有服务可用")
}

//Close 关闭
func (c *Client) Close() {
	c.managerclient.Close()
	c.wait.Wait()
	c.discovery.Close()
	c.fusing.Close()
	c.limit.Close()
}
