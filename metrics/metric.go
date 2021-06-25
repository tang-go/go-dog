package metrics

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
)

// Metric defines a metric object. Users can use it to save
// metric data. Every metric should be globally unique by name.
type Metric struct {
	Type       MetricType
	Name       string
	Help       string
	Labels     []string
	Buckets    []float64
	Objectives map[float64]float64

	vec prometheus.Collector
}

// SetGaugeValue set data for Gauge type Metric with values.
func (m *Metric) SetGaugeValue(labelValues []string, value float64) error {
	if m.Type == Untyped {
		return fmt.Errorf("metric '%s' not existed", m.Name)
	}

	if m.Type != Gauge {
		return fmt.Errorf("metric '%s' not Gauge type", m.Name)
	}
	m.vec.(*prometheus.GaugeVec).WithLabelValues(labelValues...).Set(value)
	return nil
}

// SetGaugeValue set data for Gauge type Metric with labels and values.
func (m *Metric) SetGaugeValueWithLabel(labelValues map[string]string, value float64) error {
	if m.Type == Untyped {
		return fmt.Errorf("metric '%s' not existed", m.Name)
	}

	if m.Type != Gauge {
		return fmt.Errorf("metric '%s' not Gauge type", m.Name)
	}
	m.vec.(*prometheus.GaugeVec).With(labelValues).Set(value)
	return nil
}

// Inc increases value for Counter/Gauge type metric, increments
// the counter by 1
func (m *Metric) Inc(labelValues []string) error {
	if m.Type == Untyped {
		return fmt.Errorf("metric '%s' not existed", m.Name)
	}

	if m.Type != Gauge && m.Type != Counter {
		return fmt.Errorf("metric '%s' not Gauge or Counter type", m.Name)
	}
	switch m.Type {
	case Counter:
		m.vec.(*prometheus.CounterVec).WithLabelValues(labelValues...).Inc()
		break
	case Gauge:
		m.vec.(*prometheus.GaugeVec).WithLabelValues(labelValues...).Inc()
		break
	}
	return nil
}

// Inc increases value for Counter/Gauge type metric, increments
// the counter by 1
func (m *Metric) IncWithLabel(labelValues map[string]string) error {
	if m.Type == Untyped {
		return fmt.Errorf("metric '%s' not existed", m.Name)
	}

	if m.Type != Gauge && m.Type != Counter {
		return fmt.Errorf("metric '%s' not Gauge or Counter type", m.Name)
	}
	switch m.Type {
	case Counter:
		m.vec.(*prometheus.CounterVec).With(labelValues).Inc()
		break
	case Gauge:
		m.vec.(*prometheus.GaugeVec).With(labelValues).Inc()
		break
	}
	return nil
}

// Add adds the given value to the Metric object. Only
// for Counter/Gauge type metric.
func (m *Metric) Add(labelValues []string, value float64) error {
	if m.Type == Untyped {
		return fmt.Errorf("metric '%s' not existed", m.Name)
	}

	if m.Type != Gauge && m.Type != Counter {
		return fmt.Errorf("metric '%s' not Gauge or Counter type", m.Name)
	}
	switch m.Type {
	case Counter:
		m.vec.(*prometheus.CounterVec).WithLabelValues(labelValues...).Add(value)
		break
	case Gauge:
		m.vec.(*prometheus.GaugeVec).WithLabelValues(labelValues...).Add(value)
		break
	}
	return nil
}

// Add adds the given value to the Metric object. Only
// for Counter/Gauge type metric.
func (m *Metric) AddWithLabel(labelValues map[string]string, value float64) error {
	if m.Type == Untyped {
		return fmt.Errorf("metric '%s' not existed", m.Name)
	}

	if m.Type != Gauge && m.Type != Counter {
		return fmt.Errorf("metric '%s' not Gauge or Counter type", m.Name)
	}
	switch m.Type {
	case Counter:
		m.vec.(*prometheus.CounterVec).With(labelValues).Add(value)
		break
	case Gauge:
		m.vec.(*prometheus.GaugeVec).With(labelValues).Add(value)
		break
	}
	return nil
}

// Observe is used by Histogram and Summary type metric to
// add observations.
func (m *Metric) Observe(labelValues []string, value float64) error {
	if m.Type == Untyped {
		return fmt.Errorf("metric '%s' not existed", m.Name)
	}
	if m.Type != Histogram && m.Type != Summary {
		return fmt.Errorf("metric '%s' not Histogram or Summary type", m.Name)
	}
	switch m.Type {
	case Histogram:
		m.vec.(*prometheus.HistogramVec).WithLabelValues(labelValues...).Observe(value)
		break
	case Summary:
		m.vec.(*prometheus.SummaryVec).WithLabelValues(labelValues...).Observe(value)
		break
	}
	return nil
}

// Observe is used by Histogram and Summary type metric to
// add observations.
func (m *Metric) ObserveWithLabel(labelValues map[string]string, value float64) error {
	if m.Type == Untyped {
		return fmt.Errorf("metric '%s' not existed", m.Name)
	}
	if m.Type != Histogram && m.Type != Summary {
		return fmt.Errorf("metric '%s' not Histogram or Summary type", m.Name)
	}
	switch m.Type {
	case Histogram:
		m.vec.(*prometheus.HistogramVec).With(labelValues).Observe(value)
		break
	case Summary:
		m.vec.(*prometheus.SummaryVec).With(labelValues).Observe(value)
		break
	}
	return nil
}

func (m *Metric) GetVec() prometheus.Collector {
	return m.vec
}

func GetMetric(name string) (*Metric, error) {
	return defaultMetricManager.GetMetric(name)
}

func SetGaugeValue(name string, labelValues []string, value float64) error {
	metric, err := GetMetric(name)
	if err != nil {
		return err
	}

	return metric.SetGaugeValue(labelValues, value)
}

func SetGaugeValueWithLabel(name string, labelValues map[string]string, value float64) error {
	metric, err := GetMetric(name)
	if err != nil {
		return err
	}

	return metric.SetGaugeValueWithLabel(labelValues, value)
}

func Inc(name string, labelValues []string) error {
	metric, err := GetMetric(name)
	if err != nil {
		return err
	}

	return metric.Inc(labelValues)
}

func IncWithLabel(name string, labelValues map[string]string) error {
	metric, err := GetMetric(name)
	if err != nil {
		return err
	}

	return metric.IncWithLabel(labelValues)
}

func Add(name string, labelValues []string, value float64) error {
	metric, err := GetMetric(name)
	if err != nil {
		return err
	}

	return metric.Add(labelValues, value)
}

func AddWithLabel(name string, labelValues map[string]string, value float64) error {
	metric, err := GetMetric(name)
	if err != nil {
		return err
	}

	return metric.AddWithLabel(labelValues, value)
}

func Observe(name string, labelValues []string, value float64) error {
	metric, err := GetMetric(name)
	if err != nil {
		return err
	}

	return metric.Observe(labelValues, value)
}

func ObserveWithLabel(name string, labelValues map[string]string, value float64) error {
	metric, err := GetMetric(name)
	if err != nil {
		return err
	}

	return metric.ObserveWithLabel(labelValues, value)
}
