package metrics

import (
	"fmt"
	"sync"
)

type MetricOpts struct {
	NameSpace     string         // 选填
	SystemName    string         // 必填
	MetricsValues []*MetricValue // 必填
}

// 在应用初始化的时候调用此接口
func Init(opts *MetricOpts) error {
	err := initOpts(opts)
	if err != nil {
		return err
	}
	for _, metricsValue := range defaultMetricsValues {
		if f, ok := promTypeHandler[metricsValue.ValueType]; ok {
			_, err := f(metricsValue)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("metrics init error, unknown valuetype:%s", metricsValue.ValueType)
		}
	}
	for _, metricsValue := range opts.MetricsValues {
		if f, ok := promTypeHandler[metricsValue.ValueType]; ok {
			_, err := f(metricsValue)
			if err != nil {
				return err
			}
		} else {
			return fmt.Errorf("metrics init error, unknown valuetype:%s", metricsValue.ValueType)
		}
	}
	return nil
}

func initOpts(opts *MetricOpts) error {
	defaultMetricManager = &MetricManager{
		systemName: opts.SystemName,
		metrics:    sync.Map{},
	}
	if opts.NameSpace != "" {
		defaultMetricManager.namespace = opts.NameSpace
	} else {
		defaultMetricManager.namespace = "default"
	}
	return nil
}
