package plugins

import "time"

//Context 自定义接口
type Context interface {
	Deadline() (deadline time.Time, ok bool)

	Done() <-chan struct{}

	Err() error

	Value(key interface{}) interface{}

	//SetURL 设置请求url
	SetURL(url string)

	//GetURL 获取请求url
	GetURL() string

	//SetSource 设置请求源
	SetSource(source string)

	//GetSource 获取请求源
	GetSource() string

	//GetTTL 获取超时时间
	GetTTL() int64

	//Cancel 执行取消函数
	Cancel()

	//GetTimeOut 获取超时时间
	GetTimeOut() int64

	//SetIsTest 设置是测试请求
	SetIsTest(bool)

	//GetIsTest 是否是测试请求
	GetIsTest() bool

	//SetAddress  设置请求ip
	SetAddress(address string)

	//GetAddress 获取请求ip
	GetAddress() string

	//SetTraceID  设置traceid
	SetTraceID(traceID string)

	//GetTraceID 获取traceid
	GetTraceID() string

	//SetToken 设置token
	SetToken(token string)

	//GetToken 获取token
	GetToken() string

	//SetShare 设置共享数据
	SetShare(key string, val interface{})

	//GetShare 获取全部共享数据
	GetShare() map[string]interface{}

	//GetShareByKey 获取指定共享数据
	GetShareByKey(key string) interface{}

	//SetData  设置自定义data
	SetData(key string, val interface{})

	//GetData 获取自定义data
	GetData() map[string]interface{}

	//GetDataByKey 通过key获取自定义数据
	GetDataByKey(string) interface{}

	//SetClient 设置客户端
	SetClient(cli Client)

	//GetClient 获取客户端
	GetClient() Client
}
