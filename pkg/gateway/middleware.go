package gateway

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tang-go/go-dog/metrics"
)

const (
	ReqCount      = "request_count"
	ReqDuration   = "request_duration_seconds"
	ReqSizeBytes  = "request_size_bytes"
	RespSizeBytes = "response_size_bytes"
)

const (
	Method = "method"
	Path   = "path"
	Name   = "name"
	Code   = "code"
)

var labels = []string{Name, Method, Path, Code}

func (g *Gateway) MetricMiddleware(c *gin.Context) {
	start := time.Now()
	c.Next()

	method := c.Request.Method
	path := c.Request.URL.Path
	code := strconv.Itoa(c.Writer.Status())

	labelValues := map[string]string{Name: g.name, Method: method, Path: path, Code: code}

	metric, err := metrics.GetManager().GetMetric(ReqCount)
	if err == nil && metric != nil {
		metric.IncWithLabel(labelValues)
	}

	metric, err = metrics.GetManager().GetMetric(ReqSizeBytes)
	if err == nil && metric != nil {
		metric.ObserveWithLabel(labelValues, calcRequestSize(c.Request))
	}

	metric, err = metrics.GetManager().GetMetric(RespSizeBytes)
	if err == nil && metric != nil {
		metric.ObserveWithLabel(labelValues, float64(c.Writer.Size()))
	}

	metric, err = metrics.GetManager().GetMetric(ReqDuration)
	if err == nil && metric != nil {
		metric.ObserveWithLabel(labelValues, time.Since(start).Seconds())
	}
}

func calcRequestSize(r *http.Request) float64 {
	size := 0
	if r.URL != nil {
		size = len(r.URL.String())
	}

	size += len(r.Method)
	size += len(r.Proto)

	for name, values := range r.Header {
		size += len(name)
		for _, value := range values {
			size += len(value)
		}
	}
	size += len(r.Host)

	// r.Form and r.MultipartForm are assumed to be included in r.URL.
	if r.ContentLength != -1 {
		size += int(r.ContentLength)
	}
	return float64(size)
}
