package gateway

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tang-go/go-dog/metrics"
)

const (
	ReqCount      = "http_request_count"
	ReqDuration   = "http_request_duration_seconds"
	ReqSizeBytes  = "http_request_size_bytes"
	RespSizeBytes = "http_response_size_bytes"
)

const (
	Method = "method"
	Path   = "path"
	Status = "status"
)

func metricMiddleware(c *gin.Context) {
	start := time.Now()
	c.Next()

	method := c.Request.Method
	path := c.Request.URL.Path
	status := strconv.Itoa(c.Writer.Status())

	labelValues := map[string]string{Method: method, Path: path, Status: status}

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
