package plugins

//Limit 限流插件
type Limit interface {
	//SetLimit 设置最大限制
	SetLimit(max int)

	//IsLimit 是否限制通过
	IsLimit() bool

	//Close 关闭
	Close()
}
