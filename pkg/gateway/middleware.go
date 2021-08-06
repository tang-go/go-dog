package gateway

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/tang-go/go-dog/metrics"
)

func (g *Gateway) metricMiddleware(c *gin.Context) {
	start := time.Now()
	code := strconv.Itoa(c.Writer.Status())
	url := c.Request.URL.Path
	metrics.MetricWorkingCount(g.name, url, 1)
	metrics.MetricRequestCount(g.name, url)
	c.Next()
	if code != "200" {
		metrics.MetricResponseCount(g.name, url, "false", code)
	} else {
		metrics.MetricResponseCount(g.name, url, "true", code)
	}
	metrics.MetricResponseTime(g.name, url, time.Since(start).Seconds())
	metrics.MetricWorkingCount(g.name, url, -1)
}
