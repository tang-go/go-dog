## 使用
#### 初始化处使用
1. 指标初始化
定义服务要使用的各种指标

```
var metricsValues = []*metrics.MetricValue{
  {
    ValueType: metrics.Counter,                                                                  // 提供四种指标类型Counter,Gauge,Histogram,Summary
    Name:      constant.MetricTranslateLanguage,                                                 // 本服务唯一
    Help:      "hello interface count",                                                          // 指标描述
    Labels:    []string{constant.MetricLabelSourceLanguage, constant.MetricLabelTargetLanguage}, // 标签注意顺序
  },
}
```

2. metrics初始化
在应用初始化的时候调用metrics.Init, 出入gin.Engine和上文定义的指标，此函数会执行两部分初始化。
* 为gin新增一个路径 */metrics*, 此路径可以获取用户的打点内容
* 初始化defaultMetricsManger, 并初始化应用的默认打点指标：
    * 带有method,path,status标签的http_request_count
    * 带有method,path,status标签的http_request_duration_seconds
    * 带有method,path,status标签的http_request_size_bytes
    * 带有method,path,status标签的http_response_size_bytes
    
#### 业务打点使用
1. 指标简单方法使用
```
    // 业务打点
	// metrics包维护默认的metric管理器
	// 使用IncWithLabel给带有标签的指标加1
	labels := map[string]string{constant.MetricLabelSourceLanguage: "cn", constant.MetricLabelTargetLanguage: "en"}
	err := metrics.IncWithLabel(constant.MetricTranslateLanguage, labels)
	if err != nil {
		log.Println("【Error】", err)
		c.Status(500)
		return
	}

	// 使用Inc给带有标签的指标加1，不过这个方法需要让初始化的label和value的顺序一样。如下所示：
	err = metrics.Inc(constant.MetricTranslateLanguage, []string{"cn", "en"})
	if err != nil {
		log.Println("【Error】", err)
		c.Status(500)
		return
	}
```

2. 指标复杂方法使用
```
    // 如果metrics库的方法不能完全满足你, 那么可以获取Vec后转换成相应的指标并使用
	metric, err := metrics.GetMetric(constant.MetricTranslateLanguage)
	if err != nil {
	    log.Println("【Error】", err)
	    c.Status(500)
	    return
	}
	metric.GetVec().(prometheus.Histogram).Desc()
```

## 注意
1. 服务需要定义好gin.Engine,调用 *metrics.Init* 会为服务增加一个中间件和一个路径 */metrics*
2. 定义的metrics的全名是这样拼装的 *namespace_systemname_name*
3. 全文会生成一个全局唯一的 *defaultMetricsManger*
4. 所有需要的指标需要在服务初始化的时候定义清楚
