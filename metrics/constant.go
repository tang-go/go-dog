package metrics

//指标类型
type MetricType string

const (
	//计数器
	Counter MetricType = "counter"
	//计量器
	Gauge MetricType = "gauge"
	//分布图
	Histogram MetricType = "histogram"
	//摘要
	Summary MetricType = "summary"
	//计量器 prometheus服务器发送类型信号
	Untyped MetricType = "untyped"
)
