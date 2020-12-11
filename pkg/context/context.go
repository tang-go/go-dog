package context

import (
	base "context"
	"errors"
	"time"

	"github.com/tang-go/go-dog/pkg/codec"
	"github.com/tang-go/go-dog/plugins"
)

//MyContext 自定义context
type MyContext struct {
	base.Context
	codec   plugins.Codec
	ttl     int64
	timeout int64
	traceID string
	isTest  bool
	address string
	source  string
	token   string
	url     string
	cancel  base.CancelFunc
	client  plugins.Client
	data    map[string][]byte
	share   map[string]interface{}
}

//Background 创建一个空context
func Background() plugins.Context {
	c := new(MyContext)
	c.codec = codec.NewCodec()
	c.Context = base.Background()
	c.data = make(map[string][]byte)
	c.share = make(map[string]interface{})
	return c
}

//NewContextByData 通过一个data创建
func NewContextByData(data map[string][]byte) plugins.Context {
	c := new(MyContext)
	c.codec = codec.NewCodec()
	c.Context = base.Background()
	c.data = data
	c.share = make(map[string]interface{})
	return c
}

//WithTimeout 创建一个超时context
func WithTimeout(ctx plugins.Context, ttl int64) plugins.Context {
	c := ctx.(*MyContext)
	c.ttl = ttl
	c.timeout = time.Now().UnixNano() + ttl
	newctx, cancel := base.WithTimeout(base.Background(), time.Duration(ttl))
	c.Context = newctx
	c.cancel = cancel
	return c
}

//GetTTL 获取超时时间
func (c *MyContext) GetTTL() int64 {
	return c.ttl
}

//Cancel 执行取消函数
func (c *MyContext) Cancel() {
	if c.cancel != nil {
		c.cancel()
	}
}

//GetTimeOut 获取超时时间
func (c *MyContext) GetTimeOut() int64 {
	return c.timeout
}

//SetIsTest 设置是测试请求
func (c *MyContext) SetIsTest(test bool) {
	c.isTest = test
}

//GetIsTest 是否是测试请求
func (c *MyContext) GetIsTest() bool {
	return c.isTest
}

//SetSource 设置请求源
func (c *MyContext) SetSource(source string) {
	c.source = source
}

//GetSource 获取请求源
func (c *MyContext) GetSource() string {
	return c.source
}

//SetURL 设置请求url
func (c *MyContext) SetURL(url string) {
	c.url = url
}

//GetURL 获取请求url
func (c *MyContext) GetURL() string {
	return c.url
}

//SetToken 设置token
func (c *MyContext) SetToken(token string) {
	c.token = token
}

//GetToken 获取token
func (c *MyContext) GetToken() string {
	return c.token
}

//SetAddress  设置请求ip
func (c *MyContext) SetAddress(address string) {
	c.address = address
}

//GetAddress 获取请求ip
func (c *MyContext) GetAddress() string {
	return c.address
}

//SetTraceID  设置traceid
func (c *MyContext) SetTraceID(traceID string) {
	c.traceID = traceID
}

//GetTraceID 获取traceid
func (c *MyContext) GetTraceID() string {
	return c.traceID
}

//SetShare 设置共享数据
func (c *MyContext) SetShare(key string, val interface{}) {
	c.share[key] = val
}

//GetShare 获取全部共享数据
func (c *MyContext) GetShare() map[string]interface{} {
	return c.share
}

//GetShareByKey 获取指定共享数据
func (c *MyContext) GetShareByKey(key string) interface{} {
	return c.share[key]
}

//SetData  设置自定义data
func (c *MyContext) SetData(key string, val interface{}) error {
	v, err := c.codec.EnCode("msgpack", val)
	if err != nil {
		return err
	}
	c.data[key] = v
	return nil
}

//GetData 获取自定义data
func (c *MyContext) GetData() map[string][]byte {
	return c.data
}

//GetDataByKey 获取自定义data值
func (c *MyContext) GetDataByKey(key string, val interface{}) error {
	v, ok := c.data[key]
	if ok {
		return c.codec.DeCode("msgpack", []byte(v), val)
	}
	return errors.New("no this key")
}

//SetClient 设置客户端
func (c *MyContext) SetClient(cli plugins.Client) {
	c.client = cli
}

//GetClient 获取客户端
func (c *MyContext) GetClient() plugins.Client {
	return c.client
}
