package context

import (
	base "context"
	"go-dog/plugins"
	"time"
)

//MyContext 自定义context
type MyContext struct {
	base.Context
	ttl     int64
	timeout int64
	traceID string
	isTest  bool
	address string
	cancel  base.CancelFunc
	client  plugins.Client
	data    map[string]interface{}
}

//Background 创建一个空context
func Background() plugins.Context {
	c := new(MyContext)
	c.Context = base.Background()
	c.data = make(map[string]interface{})
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

//SetData  设置自定义data
func (c *MyContext) SetData(key string, val interface{}) {
	c.data[key] = val
}

//GetData 获取自定义data
func (c *MyContext) GetData() map[string]interface{} {
	return c.data
}

//GetDataByKey 获取自定义data值
func (c *MyContext) GetDataByKey(key string) interface{} {
	return c.data[key]
}

//SetClient 设置客户端
func (c *MyContext) SetClient(cli plugins.Client) {
	c.client = cli
}

//GetClient 获取客户端
func (c *MyContext) GetClient() plugins.Client {
	return c.client
}
