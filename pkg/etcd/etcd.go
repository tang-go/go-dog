package etcd

import (
	"time"

	"github.com/coreos/etcd/clientv3"
)

var gEtcdClient *_EtcdClient

type _EtcdClient struct {
	client *clientv3.Client
}

func newEtcdClient(addr []string, timeout time.Duration) (*_EtcdClient, error) {
	conf := clientv3.Config{
		Endpoints:   addr,
		DialTimeout: timeout,
	}
	client, err := clientv3.New(conf)
	if err != nil {
		return nil, err
	}
	return &_EtcdClient{client: client}, nil
}

//InitEtcdClient 初始化客户端
func InitEtcdClient(addr []string, timeout time.Duration) (err error) {
	if gEtcdClient == nil {
		gEtcdClient, err = newEtcdClient(addr, timeout)
		if err != nil {
			return err
		}
	}
	return nil
}

//GetEtcdClient 获取etcd客户端
func GetEtcdClient() *clientv3.Client {
	if gEtcdClient != nil {
		return gEtcdClient.client
	}
	return nil
}
