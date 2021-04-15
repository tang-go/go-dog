package nacos

import (
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

//Config 配置中心
type Config struct {
	client config_client.IConfigClient
}

//NewConfig 初始化一个配置中心
func newConfig(client config_client.IConfigClient) *Config {
	c := new(Config)
	c.client = client
	return c
}

//GetConfig 获取配置
func (c *Config) GetConfig(dataID, group string) (string, error) {
	content, err := c.client.GetConfig(vo.ConfigParam{
		DataId: dataID,
		Group:  group,
	})
	if err != nil {
		return "", err
	}
	return content, nil
}

//PublishConfig 发布/修改配置
func (c *Config) PublishConfig(dataID, group, content string) (bool, error) {
	return c.client.PublishConfig(vo.ConfigParam{
		DataId:  dataID,
		Group:   group,
		Content: content,
	})
}

//ListenConfig 监听配置
func (c *Config) ListenConfig(dataID, group string, onChange func(namespace, group, dataId, data string)) error {
	return c.client.ListenConfig(vo.ConfigParam{
		DataId:   dataID,
		Group:    group,
		OnChange: onChange,
	})
}
