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

	//SetByTime 设置redis 储存键值对含有过期时间 这个函数被弃用,但是为了保证前面的代码能用,所以保留
	//param:key		键
	//param:value	值
	//param:tm		过期时间单位秒
	SetByTime(key string, value interface{}, tm int) error

	//Get 获取值
	Get(key string, value interface{}) error

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
