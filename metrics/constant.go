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

//默认系统指标
const (
	//运行服务名称
	ServiceRun = "service_run"
	//请求数
	RequestCount = "request_count"
	//请求响应数
	ResponseCount = "response_count"
	//请求响应时间 单位ms
	ResponseTime = "response_time"
	//正在工作的业务数量
	WorkingCount = "working_count"
	//请求byte总大小
	RequestBytes = "request_bytes"
	//响应byte总大小
	ResponseBytes = "response_bytes"
)

//默认label
const (
	Method  = "method"
	Name    = "name"
	Success = "success"
	Code    = "code"
)

//注册默认指标
var defaultMetricsValues = []*MetricValue{
	{
		ValueType: Gauge,
		Name:      ServiceRun,
		Help:      "Counter. total run service count",
		Labels:    []string{Name},
	},
	{
		ValueType: Counter,
		Name:      RequestCount,
		Help:      "Counter. total request count",
		Labels:    []string{Name, Method},
	}, {
		ValueType: Counter,
		Name:      ResponseCount,
		Help:      "Counter. total response count",
		Labels:    []string{Name, Method, Success, Code},
	}, {
		ValueType: Histogram,
		Name:      ResponseTime,
		Help:      "Histogram. request latencies in seconds",
		Labels:    []string{Name, Method},
	}, {
		ValueType: Gauge,
		Name:      WorkingCount,
		Help:      "Gauge. total working ncount",
		Labels:    []string{Name, Method},
	},
	{
		ValueType: Summary,
		Name:      RequestBytes,
		Help:      "Summary. total request bytes size",
		Labels:    []string{Name, Method},
	},
	{
		ValueType: Summary,
		Name:      ResponseBytes,
		Help:      "Summary. total response bytes size",
		Labels:    []string{Name, Method},
	},
}

//MetricResponseBytes 响应时间指标
func MetricResponseBytes(name, method string, size float64) {
	metric, err := GetManager().GetMetric(ResponseBytes)
	if err == nil && metric != nil {
		metric.ObserveWithLabel(map[string]string{Name: name, Method: method}, size)
	}
}

//MetricRequestBytes 请求时间指标
func MetricRequestBytes(name, method string, size float64) {
	metric, err := GetManager().GetMetric(RequestBytes)
	if err == nil && metric != nil {
		metric.ObserveWithLabel(map[string]string{Name: name, Method: method}, size)
	}
}

//MetricWorkingCount 正在运行的任务指标
func MetricWorkingCount(name, method string, count float64) {
	metric, err := GetManager().GetMetric(WorkingCount)
	if err == nil && metric != nil {
		metric.AddWithLabel(map[string]string{Name: name, Method: method}, count)
	}
}

//MetricResponseTime 响应时间指标
func MetricResponseTime(name, method string, seconds float64) {
	metric, err := GetManager().GetMetric(ResponseTime)
	if err == nil && metric != nil {
		metric.ObserveWithLabel(map[string]string{Name: name, Method: method}, seconds)
	}
}

//MetricResponseCount 响应数指标
func MetricResponseCount(name, method, success, code string) {
	metric, err := GetManager().GetMetric(ResponseCount)
	if err == nil && metric != nil {
		metric.IncWithLabel(map[string]string{Name: name, Method: method, Success: success, Code: code})
	}
}

//MetricRequestCount 请求数指标
func MetricRequestCount(name, method string) {
	metric, err := GetManager().GetMetric(RequestCount)
	if err == nil && metric != nil {
		metric.IncWithLabel(map[string]string{Name: name, Method: method})
	}
}

//MetricServiceRun 运行服务指标
func MetricServiceRun(name string, count float64) {
	metric, err := GetManager().GetMetric(ServiceRun)
	if err == nil && metric != nil {
		metric.AddWithLabel(map[string]string{Name: name}, count)
	}
}
