package metrics

type MetricValue struct {
	ValueType MetricType
	Name      string
	Help      string
	Labels    []string
}

var (
	promTypeHandler = map[MetricType]func(metric *MetricValue) (*Metric, error){
		Counter:   GenerateCounter,
		Gauge:     GenerateGauge,
		Histogram: GenerateHistogram,
		Summary:   GenerateSummary,
	}
)

func GenerateCounter(metricValue *MetricValue) (*Metric, error) {
	metric, err := defaultMetricManager.GenerateCounter(metricValue)
	if err != nil {
		return nil, err
	}

	return metric, nil
}

func GenerateGauge(metricValue *MetricValue) (*Metric, error) {
	metric, err := defaultMetricManager.GenerateGauge(metricValue)
	if err != nil {
		return nil, err
	}

	return metric, nil
}

func GenerateHistogram(metricValue *MetricValue) (*Metric, error) {
	metric, err := defaultMetricManager.GenerateHistogram(metricValue)
	if err != nil {
		return nil, err
	}

	return metric, nil
}

func GenerateSummary(metricValue *MetricValue) (*Metric, error) {
	metric, err := defaultMetricManager.GenerateSummary(metricValue)
	if err != nil {
		return nil, err
	}

	return metric, nil
}
