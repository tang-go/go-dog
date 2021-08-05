package gateway

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tang-go/go-dog/metrics"
)

const (
	RequestTpsCount = "request_tps_count"
	RequestQpsCount = "request_qps_count"
	RequestSeconds  = "request_seconds"
)

const (
	Method  = "method"
	Name    = "name"
	Success = "success"
	Code    = "code"
)

func (g *Gateway) metricMiddleware(c *gin.Context) {
	metric, err := metrics.GetManager().GetMetric(RequestQpsCount)
	if err == nil && metric != nil {
		metric.IncWithLabel(map[string]string{Name: g.name, Method: c.Request.URL.Path})
	}
	start := time.Now()
	c.Next()
	code := strconv.Itoa(c.Writer.Status())
	url := c.Request.URL.Path
	labelValues := map[string]string{Name: g.name, Method: url, Success: "true", Code: code}
	if code != "200" {
		labelValues[Success] = "false"
	}
	metric, err = metrics.GetManager().GetMetric(RequestTpsCount)
	if err == nil && metric != nil {
		metric.IncWithLabel(labelValues)
	}
	metric, err = metrics.GetManager().GetMetric(RequestSeconds)
	if err == nil && metric != nil {
		metric.ObserveWithLabel(labelValues, time.Since(start).Seconds())
	}
}
