package redis

import (
	"encoding/json"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

//Redis 单机Redis
type Redis struct {
	addrs  []string
	addr   string
	pass   string
	mem    *Mem
	client *redis.Client
}

//CreateOne 创建单机redis
func CreateOne(ip string, p string) (*Redis, error) {
	redisCluster := new(Redis)
	redisCluster.addr = ip
	redisCluster.pass = p
	redisCluster.mem = NewMem()
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

//Close 关闭
func (pointer *Redis) Close() {
	pointer.mem.Close()
	pointer.client.Close()
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
func (pointer *Redis) Del(key string) (int, error) {
	if pointer.client == nil {
		//进行重连
		pointer.funcConnect()
	}
	//内存删除
	pointer.mem.Del(key)
	return 0, pointer.client.Del(key).Err()
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
func (pointer *Redis) SetByTime(key string, value interface{}, tm int64) error {
	if pointer.client == nil {
		//进行重连
		pointer.funcConnect()
	}
	buff, _ := json.Marshal(value)
	//内存存放
	pointer.mem.Set(key, string(buff), tm)
	return pointer.client.Set(key, string(buff), time.Second*time.Duration(tm)).Err()
}

//Get 获取
/*
*设置redis 根据键获取值
*参数说明:
*@param:key		键
 */
func (pointer *Redis) Get(key string, value interface{}) error {
	v, ok := pointer.mem.Get(key)
	if !ok {
		//内存里面没有则redis获取
		if pointer.client == nil {
			//进行重连
			pointer.funcConnect()
		}
		redisvalue, err := pointer.client.Get(key).Result()
		if err != nil {
			return err
		}
		ttl, err := pointer.client.TTL(key).Result()
		if err == nil {
			//没有超时进行同步
			if ttl > 0 {
				pointer.mem.Set(key, redisvalue, int64(ttl))
			}
		}
		v = redisvalue
	}
	err := json.Unmarshal([]byte(v), value)
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

//Zadd  有序集合
func (pointer *Redis) Zadd(key string, score int64, value string) (int64, error) {
	if pointer.client == nil {
		//进行重连
		pointer.funcConnect()
	}
	return pointer.client.ZAdd(key, redis.Z{Score: float64(score), Member: value}).Result()
}

//ZRange  获取有序集合成员(升序)
func (pointer *Redis) ZRange(key string, start int64, stop int64) ([]string, error) {
	if pointer.client == nil {
		//进行重连
		pointer.funcConnect()
	}
	return pointer.client.ZRange(key, start, stop).Result()
}

//ZRevRangeByScore  更具积分降序获取成员
func (pointer *Redis) ZRevRangeByScore(key string, max int64, min int64) ([]string, error) {
	if pointer.client == nil {
		//进行重连
		pointer.funcConnect()
	}
	return pointer.client.ZRevRangeByScore(key, redis.ZRangeBy{Max: strconv.FormatInt(max, 10), Min: strconv.FormatInt(min, 10)}).Result()
}

//ZRemRangeByScore  删除指定积分内的成员
func (pointer *Redis) ZRemRangeByScore(key string, min int64, max int64) (int64, error) {
	if pointer.client == nil {
		//进行重连
		pointer.funcConnect()
	}
	return pointer.client.ZRemRangeByScore(key, strconv.FormatInt(min, 10), strconv.FormatInt(max, 10)).Result()
}

//ZRem  删除指定会员
func (pointer *Redis) ZRem(key string, value string) (int64, error) {
	if pointer.client == nil {
		//进行重连
		pointer.funcConnect()
	}
	return pointer.client.ZRem(key, value).Result()
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
