package etcd

import (
	"time"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/concurrency"
	"github.com/tang-go/go-dog/plugins"
)

//Etcd Etcd分布式锁
type Etcd struct {
	cli *clientv3.Client
}

//NewEtcd 创建一个etcd客户端
func NewEtcd(cfg plugins.Cfg) *Etcd {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   cfg.GetEtcd(),
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(err.Error())
	}
	return &Etcd{
		cli: cli,
	}
}

//GetMutex 获取分布式锁对象
func (e *Etcd) GetMutex(key string) (*concurrency.Mutex, error) {
	session, err := concurrency.NewSession(e.cli)
	if err != nil {
		return nil, err
	}
	return concurrency.NewMutex(session, key), nil
}

//Close 关闭
func (e *Etcd) Close() {
	if e.cli != nil {
		e.cli.Close()
	}
}
