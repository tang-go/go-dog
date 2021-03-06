package register

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/tang-go/go-dog/serviceinfo"

	"github.com/coreos/etcd/clientv3"
)

//EtcdRegister Etcd 服务注册
type EtcdRegister struct {
	client *clientv3.Client //etcd 客户端
	rpcID  clientv3.LeaseID //rpc服务接口注册id
	apiID  clientv3.LeaseID //api服务注册接口id
	ttl    int64            //时间
}

//NewEtcdRegister 初始化一个etcd服务注册中心
func NewEtcdRegister(address []string) *EtcdRegister {
	conf := clientv3.Config{
		Endpoints:   address,
		DialTimeout: time.Duration(2) * time.Second,
	}
	client, err := clientv3.New(conf)
	if err != nil {
		panic(err.Error())
	}
	return &EtcdRegister{
		ttl:    2,
		client: client,
	}
}

//RegisterRPCService 注册RPC服务
func (s *EtcdRegister) RegisterRPCService(ctx context.Context, info *serviceinfo.ServiceInfo) error {
	key := "rpc/" + fmt.Sprintf("%s:%d", info.Address, info.Port)
	info.Key = key
	val, _ := json.Marshal(info)

	//设置租约时间
	resp, err := s.client.Grant(ctx, s.ttl)
	if err != nil {
		panic(err)
	}
	//注册服务并绑定租约
	_, err = s.client.Put(ctx, key, string(val), clientv3.WithLease(resp.ID))
	if err != nil {
		panic(err)
	}
	//设置续租 定期发送需求请求
	leaseRespChan, err := s.client.KeepAlive(ctx, resp.ID)
	if err != nil {
		panic(err)
	}
	go func() {
		for range leaseRespChan {
		}
	}()
	s.rpcID = resp.ID
	return nil
}

//RegisterAPIService 注册API服务
func (s *EtcdRegister) RegisterAPIService(ctx context.Context, info *serviceinfo.ServiceInfo) error {
	key := "api/" + fmt.Sprintf("%s:%d", info.Address, info.Port)
	info.Key = key
	val, _ := json.Marshal(info)

	//设置租约时间
	resp, err := s.client.Grant(ctx, s.ttl)
	if err != nil {
		panic(err)
	}
	//注册服务并绑定租约
	_, err = s.client.Put(ctx, key, string(val), clientv3.WithLease(resp.ID))
	if err != nil {
		panic(err)
	}
	//设置续租 定期发送需求请求
	leaseRespChan, err := s.client.KeepAlive(ctx, resp.ID)
	if err != nil {
		panic(err)
	}
	go func() {
		for range leaseRespChan {
		}
	}()
	s.apiID = resp.ID
	return nil
}

// Cancellation 注销服务
func (s *EtcdRegister) Cancellation() error {
	//撤销api接口注册
	if s.apiID > 0 {
		if _, err := s.client.Revoke(context.Background(), s.apiID); err != nil {
			return err
		}
	}
	//撤销rpc接口注册
	if s.rpcID > 0 {
		if _, err := s.client.Revoke(context.Background(), s.rpcID); err != nil {
			return err
		}
	}
	return s.client.Close()
}
