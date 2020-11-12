package redis

import (
	"encoding/json"
	"time"

	"github.com/go-redis/redis"
)

//Redis 单机Redis
type Redis struct {
	addrs  []string
	addr   string
	pass   string
	client *redis.Client
}

//CreateOne 创建单机redis
func CreateOne(ip string, p string) (*Redis, error) {
	redisCluster := new(Redis)
	redisCluster.addr = ip
	redisCluster.pass = p
	return redisCluster, redisCluster.funcConnect()
}

func (pointer *Redis) funcConnect() error {
	pointer.client = redis.NewClient(&redis.Options{
		Addr:     pointer.addr,
		Password: pointer.pass,
	})
	_, err := pointer.client.Ping().Result()
	if err != nil {
		return err
	}
	return nil
}

//Set 设置
/*
*设置redis 储存键值对含有过期时间 这个函数被弃用,但是为了保证前面的代码能用,所以保留
*参数说明:
*@param:key		键
*@param:value	值
 */
func (pointer *Redis) Set(key string, value interface{}) error {
	if pointer.client == nil {
		//进行重连
		pointer.funcConnect()
	}
	err := pointer.client.Set(key, value, 0).Err()
	return err
}

//Del 删除
func (pointer *Redis) Del(Key string) (int, error) {
	if pointer.client == nil {
		//进行重连
		pointer.funcConnect()
	}
	return 0, pointer.client.Del(Key).Err()
}

//FlushAll 清空
/*
*清空所有数据库
*参数说明:
*@param:key		键
*@param:value	值
 */
func (pointer *Redis) FlushAll() error {
	if pointer.client == nil {
		//进行重连
		pointer.funcConnect()
	}
	return pointer.client.FlushAll().Err()
}

//Ping 验活
/*
*测试redis是否能够连通
*参数说明:无
 */
func (pointer *Redis) Ping() (string, error) {
	if pointer.client == nil {
		//进行重连
		pointer.funcConnect()
	}
	return pointer.client.Ping().Result()
}

//SetByTime 设置
/*
*设置redis 储存键值对含有过期时间 这个函数被弃用,但是为了保证前面的代码能用,所以保留
*参数说明:
*@param:key		键
*@param:value		值
*@param:tm		过期时间单位秒
 */
func (pointer *Redis) SetByTime(key string, value interface{}, tm int) error {
	if pointer.client == nil {
		//进行重连
		pointer.funcConnect()
	}
	buff, _ := json.Marshal(value)
	return pointer.client.Set(key, string(buff), time.Second*time.Duration(tm)).Err()
}

//Get 获取
/*
*设置redis 根据键获取值
*参数说明:
*@param:key		键
 */
func (pointer *Redis) Get(key string, value interface{}) error {
	if pointer.client == nil {
		//进行重连
		pointer.funcConnect()
	}
	v, err := pointer.client.Get(key).Result()
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(v), value)
	return err
}

//SetNx 设置nx
/*
* 如果不存在相关的key,value 则设置,否则不设置
* 参数说明:
* @param:key   redis中的key
* @param:value redis中的value
* @param:tm 	redis中的超时
 */
func (pointer *Redis) SetNx(key string, value interface{}, tm int) (bool, error) {
	if pointer.client == nil {
		//进行重连
		pointer.funcConnect()
	}
	return pointer.client.SetNX(key, value, time.Second*time.Duration(tm)).Result()
}

//GetSet  获取设置
func (pointer *Redis) GetSet(key string, value interface{}) (string, error) {
	if pointer.client == nil {
		//进行重连
		pointer.funcConnect()
	}
	return pointer.client.GetSet(key, value).Result()
}

//IncrBy  递增
func (pointer *Redis) IncrBy(key string, value int64) (int64, error) {
	if pointer.client == nil {
		//进行重连
		pointer.funcConnect()
	}
	return pointer.client.IncrBy(key, value).Result()
}

//Sadd  集合
func (pointer *Redis) Sadd(key string, value string) (int64, error) {
	if pointer.client == nil {
		//进行重连
		pointer.funcConnect()
	}
	return pointer.client.SAdd(key, value).Result()
}

//SCard  获取集合成员数
func (pointer *Redis) SCard(key string) (int64, error) {
	if pointer.client == nil {
		//进行重连
		pointer.funcConnect()
	}
	return pointer.client.SCard(key).Result()
}

//SRem  删除集合成员数
func (pointer *Redis) SRem(key string, member string) (int64, error) {
	if pointer.client == nil {
		//进行重连
		pointer.funcConnect()
	}
	return pointer.client.SRem(key, member).Result()
}

//SMembers  获取集合
func (pointer *Redis) SMembers(key string) (r []string, e error) {
	if pointer.client == nil {
		//进行重连
		pointer.funcConnect()
	}
	array, err := pointer.client.SMembers(key).Result()
	if err != nil {
		return nil, err
	}
	for _, v := range array {
		r = append(r, v)
	}
	return
}

//LPush 添加列表
func (pointer *Redis) LPush(key string, vali interface{}) (int64, error) {
	if pointer.client == nil {
		//进行重连
		pointer.funcConnect()
	}
	return pointer.client.LPush(key, vali).Result()
}

//LRange 遍历列表
func (pointer *Redis) LRange(key string, start, stop int64) ([]string, error) {
	if pointer.client == nil {
		//进行重连
		pointer.funcConnect()
	}
	return pointer.client.LRange(key, start, stop).Result()
}
