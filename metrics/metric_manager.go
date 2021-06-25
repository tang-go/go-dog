package metrics

import (
	"errors"
	"fmt"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type MetricManager struct {
	namespace  string
	systemName string
	metrics    sync.Map
}

var defaultMetricManager *MetricManager

//GetManager Get Manager
func GetManager() *MetricManager {
	return defaultMetricManager
}

// Generate Counter metrics
func (p *MetricManager) GenerateCounter(metricValue *MetricValue) (*Metric, error) {
	if metricValue.Name == "" {
		return nil, errors.New("metric name is null")
	}

	vec := promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: defaultMetricManager.namespace,
		Subsystem: defaultMetricManager.systemName,
		Name:      metricValue.Name,
		Help:      metricValue.Help,
	}, metricValue.Labels)
	metric := &Metric{
		Type:   metricValue.ValueType,
		Name:   metricValue.Name,
		Help:   metricValue.Help,
		Labels: metricValue.Labels,
		vec:    vec,
	}

	_, ok := p.metrics.LoadOrStore(metricValue.Name, metric)
	if ok {
		return nil, fmt.Errorf("metrics alreay has metric:%s", metricValue.Name)
	}

	return metric, nil
}

// Generate Gauge metric
func (p *MetricManager) GenerateGauge(metricValue *MetricValue) (*Metric, error) {
	if metricValue.Name == "" {
		return nil, errors.New("metric name is null")
	}

	vec := promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: defaultMetricManager.namespace,
		Subsystem: defaultMetricManager.systemName,
		Name:      metricValue.Name,
		Help:      metricValue.Help,
	}, metricValue.Labels)
	metric := &Metric{
		Type:   metricValue.ValueType,
		Name:   metricValue.Name,
		Help:   metricValue.Help,
		Labels: metricValue.Labels,
		vec:    vec,
	}

	_, ok := p.metrics.LoadOrStore(metricValue.Name, metric)
	if ok {
		return nil, fmt.Errorf("metrics alreay has metric:%s", metricValue.Name)
	}

	return metric, nil
}

// Generate Histogram metric
func (p *MetricManager) GenerateHistogram(metricValue *MetricValue) (*Metric, error) {
	if metricValue.Name == "" {
		return nil, errors.New("metric name is null")
	}

	vec := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: defaultMetricManager.namespace,
		Subsystem: defaultMetricManager.systemName,
		Name:      metricValue.Name,
		Help:      metricValue.Help,
	}, metricValue.Labels)
	metric := &Metric{
		Type:   metricValue.ValueType,
		Name:   metricValue.Name,
		Help:   metricValue.Help,
		Labels: metricValue.Labels,
		vec:    vec,
	}

	_, ok := p.metrics.LoadOrStore(metricValue.Name, metric)
	if ok {
		return nil, fmt.Errorf("metrics alreay has metric:%s", metricValue.Name)
	}

	return metric, nil
}

// Generate Summary metric
func (p *MetricManager) GenerateSummary(metricValue *MetricValue) (*Metric, error) {
	if metricValue.Name == "" {
		return nil, errors.New("metric name is null")
	}

	vec := promauto.NewSummaryVec(prometheus.SummaryOpts{
		Namespace: defaultMetricManager.namespace,
		Subsystem: defaultMetricManager.systemName,
		Name:      metricValue.Name,
		Help:      metricValue.Help,
	}, metricValue.Labels)
	metric := &Metric{
		Type:   metricValue.ValueType,
		Name:   metricValue.Name,
		Help:   metricValue.Help,
		Labels: metricValue.Labels,
		vec:    vec,
	}

	_, ok := p.metrics.LoadOrStore(metricValue.Name, metric)
	if ok {
		return nil, fmt.Errorf("metrics alreay has metric:%s", metricValue.Name)
	}

	return metric, nil
}

func (p *MetricManager) GetMetric(name string) (*Metric, error) {
	metric, ok := p.metrics.Load(name)
	if !ok {
		return nil, fmt.Errorf("dont contain metric(%s)", name)
	}
	return metric.(*Metric), nil
}
