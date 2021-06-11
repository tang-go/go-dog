package discovery

import (
	"context"
	"encoding/json"
	"time"

	"github.com/tang-go/go-dog/log"
	"github.com/tang-go/go-dog/serviceinfo"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
)

//EtcdDiscovery 服务发现
type EtcdDiscovery struct {
	client                 *clientv3.Client //etcd 客户端
	rpcServcieOnlineNotice func(string, *serviceinfo.ServiceInfo)
	rpcServcieOffineNotice func(string)
	apiServcieOnlineNotice func(string, *serviceinfo.ServiceInfo)
	apiServcieOffineNotice func(string)
}

//NewEtcdDiscovery  新建发现服务
func NewEtcdDiscovery(address []string) *EtcdDiscovery {
	conf := clientv3.Config{
		Endpoints:   address,
		DialTimeout: time.Duration(2) * time.Second,
	}
	client, err := clientv3.New(conf)
	if err != nil {
		panic(err.Error())
	}
	return &EtcdDiscovery{
		client: client,
	}
}

//RegRPCServiceOnlineNotice 注册RPC服务上线通知
func (d *EtcdDiscovery) RegRPCServiceOnlineNotice(f func(string, *serviceinfo.ServiceInfo)) {
	d.rpcServcieOnlineNotice = f
}

//RegRPCServiceOfflineNotice 注册RPC服务下线通知
func (d *EtcdDiscovery) RegRPCServiceOfflineNotice(f func(string)) {
	d.rpcServcieOffineNotice = f
}

//RegAPIServiceOnlineNotice 注册API服务上线通知
func (d *EtcdDiscovery) RegAPIServiceOnlineNotice(f func(string, *serviceinfo.ServiceInfo)) {
	d.apiServcieOnlineNotice = f
}

//RegAPIServiceOfflineNotice 注册API服务下线通知
func (d *EtcdDiscovery) RegAPIServiceOfflineNotice(f func(string)) {
	d.apiServcieOffineNotice = f
}

//WatchRPCService 开始RPC服务发现
func (d *EtcdDiscovery) WatchRPCService() {
	//根据前缀获取现有的key
	resp, err := d.client.Get(context.Background(), "rpc/", clientv3.WithPrefix())
	if err != nil {
		panic(err)
	}
	for _, ev := range resp.Kvs {
		if d.rpcServcieOnlineNotice != nil {
			info := serviceinfo.ServiceInfo{}
			if err := json.Unmarshal(ev.Value, &info); err != nil {
				continue
			}
			d.rpcServcieOnlineNotice(string(ev.Key), &info)
		}
	}
	go func() {
		rch := d.client.Watch(context.Background(), "rpc/", clientv3.WithPrefix())
		for wresp := range rch {
			for _, ev := range wresp.Events {
				switch ev.Type {
				case mvccpb.PUT: //修改或者新增
					if d.rpcServcieOnlineNotice != nil {
						info := serviceinfo.ServiceInfo{}
						if err := json.Unmarshal(ev.Kv.Value, &info); err != nil {
							log.Errorln(err.Error(), ev.Kv.Key, ev.Kv.Value)
							return
						}
						d.rpcServcieOnlineNotice(string(ev.Kv.Key), &info)
					}
				case mvccpb.DELETE: //删除
					if d.rpcServcieOffineNotice != nil {
						d.rpcServcieOffineNotice(string(ev.Kv.Key))
					}
				}
			}
		}
	}()
}

//WatchAPIService 开始API服务发现
func (d *EtcdDiscovery) WatchAPIService() {
	//根据前缀获取现有的key
	resp, err := d.client.Get(context.Background(), "api/", clientv3.WithPrefix())
	if err != nil {
		panic(err)
	}
	for _, ev := range resp.Kvs {
		if d.apiServcieOnlineNotice != nil {
			info := serviceinfo.ServiceInfo{}
			if err := json.Unmarshal(ev.Value, &info); err != nil {
				continue
			}
			d.apiServcieOnlineNotice(string(ev.Key), &info)
		}
	}
	go func() {
		rch := d.client.Watch(context.Background(), "api/", clientv3.WithPrefix())
		for wresp := range rch {
			for _, ev := range wresp.Events {
				switch ev.Type {
				case mvccpb.PUT: //修改或者新增
					if d.apiServcieOnlineNotice != nil {
						info := serviceinfo.ServiceInfo{}
						if err := json.Unmarshal(ev.Kv.Value, &info); err != nil {
							log.Errorln(err.Error(), ev.Kv.Key, ev.Kv.Value)
							return
						}
						d.apiServcieOnlineNotice(string(ev.Kv.Key), &info)
					}
				case mvccpb.DELETE: //删除
					if d.apiServcieOffineNotice != nil {
						d.apiServcieOffineNotice(string(ev.Kv.Key))
					}
				}
			}
		}
	}()
}

//Close 关闭服务
func (d *EtcdDiscovery) Close() error {
	return d.client.Close()
}
