package redis

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

//Cluster redis集群
type Cluster struct {
	addrs   []string
	addr    string
	pass    string
	clients *redis.ClusterClient
}

//CreateCluster 创建集群
func CreateCluster(ip []string, p string) (*Cluster, error) {
	redisCluster := new(Cluster)
	redisCluster.addrs = ip
	redisCluster.pass = p
	return redisCluster, redisCluster.funcConnect()
}

func (pointer *Cluster) funcConnect() error {
	pointer.clients = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    pointer.addrs, //set redis cluster url
		Password: pointer.pass,  //set password
	})

	_, err := pointer.clients.Ping().Result()
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
func (pointer *Cluster) Set(key string, value interface{}) error {
	if pointer.clients == nil {
		//进行重连
		pointer.funcConnect()
	}
	err := pointer.clients.Set(key, value, 0).Err()
	return err
}

//Del 删除
func (pointer *Cluster) Del(Key string) (int, error) {
	if pointer.clients == nil {
		//进行重连
		pointer.funcConnect()
	}
	return 0, pointer.clients.Del(Key).Err()
}

//FlushAll 清空所有数据库
/*
*清空所有数据库
*参数说明:
*@param:key		键
*@param:value	值
 */
func (pointer *Cluster) FlushAll() error {
	if pointer.clients == nil {
		//进行重连
		pointer.funcConnect()
	}
	return pointer.clients.FlushAll().Err()
}

//Ping 验活
/*
*测试redis是否能够连通
*参数说明:无
 */
func (pointer *Cluster) Ping() (string, error) {
	if pointer.clients == nil {
		//进行重连
		pointer.funcConnect()
	}
	return pointer.clients.Ping().Result()
}

//SetByTime 设置超时
/*
*设置redis 储存键值对含有过期时间 这个函数被弃用,但是为了保证前面的代码能用,所以保留
*参数说明:
*@param:key		键
*@param:value		值
*@param:tm		过期时间单位秒
 */
func (pointer *Cluster) SetByTime(key string, value interface{}, tm int) error {
	if pointer.clients == nil {
		//进行重连
		pointer.funcConnect()
	}
	buff, _ := json.Marshal(value)
	return pointer.clients.Set(key, string(buff), time.Second*time.Duration(tm)).Err()
}

//Get 获取
/*
*设置redis 根据键获取值
*参数说明:
*@param:key		键
 */
func (pointer *Cluster) Get(key string, value interface{}) error {
	if pointer.clients == nil {
		//进行重连
		pointer.funcConnect()
	}
	v, err := pointer.clients.Get(key).Result()
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(v), value)
	return err
}

//SetNx 设置
/*
* 如果不存在相关的key,value 则设置,否则不设置
* 参数说明:
* @param:key   redis中的key
* @param:value redis中的value
* @param:tm 	redis中的超时
 */
func (pointer *Cluster) SetNx(key string, value interface{}, tm int) (bool, error) {
	if pointer.clients == nil {
		//进行重连
		pointer.funcConnect()
	}
	return pointer.clients.SetNX(key, value, time.Second*time.Duration(tm)).Result()
}

//GetSet 获取设置
func (pointer *Cluster) GetSet(key string, value interface{}) (string, error) {
	if pointer.clients == nil {
		//进行重连
		pointer.funcConnect()
	}
	return pointer.clients.GetSet(key, value).Result()
}

//IncrBy  递增
func (pointer *Cluster) IncrBy(key string, value int64) (int64, error) {
	if pointer.clients == nil {
		//进行重连
		pointer.funcConnect()
	}
	return pointer.clients.IncrBy(key, value).Result()
}

//Zadd  有序集合
func (pointer *Cluster) Zadd(key string, score int64, value string) (int64, error) {
	if pointer.clients == nil {
		//进行重连
		pointer.funcConnect()
	}
	return pointer.clients.ZAdd(key, redis.Z{Score: float64(score), Member: value}).Result()
}

//ZRange  获取有序集合成员(升序)
func (pointer *Cluster) ZRange(key string, start int64, stop int64) ([]string, error) {
	if pointer.clients == nil {
		//进行重连
		pointer.funcConnect()
	}
	return pointer.clients.ZRange(key, start, stop).Result()
}

//ZRevRangeByScore  更具积分降序获取成员
func (pointer *Cluster) ZRevRangeByScore(key string, max int64, min int64) ([]string, error) {
	if pointer.clients == nil {
		//进行重连
		pointer.funcConnect()
	}
	return pointer.clients.ZRevRangeByScore(key, redis.ZRangeBy{Max: strconv.FormatInt(max, 10), Min: strconv.FormatInt(min, 10)}).Result()
}

//ZRemRangeByScore  删除指定积分内的成员
func (pointer *Cluster) ZRemRangeByScore(key string, min int64, max int64) (int64, error) {
	if pointer.clients == nil {
		//进行重连
		pointer.funcConnect()
	}
	return pointer.clients.ZRemRangeByScore(key, strconv.FormatInt(min, 10), strconv.FormatInt(max, 10)).Result()
}

//ZRem  删除指定会员
func (pointer *Cluster) ZRem(key string, value string) (int64, error) {
	if pointer.clients == nil {
		//进行重连
		pointer.funcConnect()
	}
	return pointer.clients.ZRem(key, value).Result()
}

//Sadd  集合
func (pointer *Cluster) Sadd(key string, value string) (int64, error) {
	if pointer.clients == nil {
		//进行重连
		pointer.funcConnect()
	}
	return pointer.clients.SAdd(key, value).Result()
}

//SCard  获取集合成员数
func (pointer *Cluster) SCard(key string) (int64, error) {
	if pointer.clients == nil {
		//进行重连
		pointer.funcConnect()
	}
	return pointer.clients.SCard(key).Result()
}

//SRem  删除集合成员数
func (pointer *Cluster) SRem(key string, member string) (int64, error) {
	if pointer.clients == nil {
		//进行重连
		pointer.funcConnect()
	}
	return pointer.clients.SRem(key, member).Result()
}

//SMembers  获取集合
func (pointer *Cluster) SMembers(key string) (r []string, e error) {
	if pointer.clients == nil {
		//进行重连
		pointer.funcConnect()
	}
	array, err := pointer.clients.SMembers(key).Result()
	if err != nil {
		return nil, err
	}
	for _, v := range array {
		r = append(r, v)
	}
	return
}

//LPush 添加列表
func (pointer *Cluster) LPush(key string, vali interface{}) (int64, error) {
	if pointer.clients == nil {
		//进行重连
		pointer.funcConnect()
	}
	return pointer.clients.LPush(key, vali).Result()
}

//LRange 遍历列表
func (pointer *Cluster) LRange(key string, start, stop int64) ([]string, error) {
	if pointer.clients == nil {
		//进行重连
		pointer.funcConnect()
	}
	return pointer.clients.LRange(key, start, stop).Result()
}
