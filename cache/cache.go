package cache

import (
	"github.com/tang-go/go-dog/cache/redis"
	"github.com/tang-go/go-dog/plugins"
)

//Inter 接口
type Inter interface {

	//FlushAll 清空所有数据库
	FlushAll() error

	//Del 删除key
	Del(Key string) (int, error)

	//SetByTime 设置redis 过期时间为单位秒
	SetByTime(key string, value interface{}, tm int) error

	//Get 获取值
	Get(key string, value interface{}) error

	//Zadd  有序集合
	Zadd(key string, score int64, value string) (int64, error)

	//ZRange  获取有序集合成员(升序)
	ZRange(key string, start int64, stop int64) ([]string, error)

	//ZRevRangeByScore  更具积分降序获取成员
	ZRevRangeByScore(key string, max int64, min int64) ([]string, error)

	//ZRemRangeByScore  删除指定积分内的成员
	ZRemRangeByScore(key string, min int64, max int64) (int64, error)

	//Sadd  集合
	Sadd(key string, value string) (int64, error)

	//SCard  获取集合成员数
	SCard(key string) (int64, error)

	//SRem  删除集合成员数
	SRem(key string, member string) (int64, error)

	//SMembers  获取集合
	SMembers(key string) (r []string, e error)
}

//Cache 缓存
type Cache struct {
	cache Inter
}

//NewCache 创建缓存
func NewCache(cfg plugins.Cfg) *Cache {
	c := new(Cache)
	address := cfg.GetRedis()
	if len(address) > 1 {
		cache, err := redis.CreateCluster(address, "")
		if err != nil {
			panic(err.Error())
		}
		c.cache = cache
	}
	if len(address) == 1 {
		cache, err := redis.CreateOne(address[0], "")
		if err != nil {
			panic(err.Error())
		}
		c.cache = cache
	}
	return c
}

//GetCache 获取缓存
func (c *Cache) GetCache() Inter {
	return c.cache
}
