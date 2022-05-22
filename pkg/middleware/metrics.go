package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/go-eagle/eagle/pkg/metric"
)

var namespace = "eagle"

var (
	labels = []string{"status", "endpoint", "method", "service"}

	// QPS
	reqCount = metric.NewCounterVec(
		&metric.CounterVecOpts{
			Namespace: namespace,
			Name:      "http_request_count_total",
			Help:      "Total number of HTTP requests made.",
			Labels:    labels,
		})

	// The QPS currently processing the request
	curReqCount = metric.NewGaugeVec(
		&metric.GaugeVecOpts{
			Namespace: namespace,
			Name:      "http_request_in_flight",
			Help:      "Current number of http requests in flight.",
			Labels:    labels,
		})

	// interface response time
	reqDuration = metric.NewHistogramVec(
		&metric.HistogramVecOpts{
			Namespace: namespace,
			Name:      "http_request_duration_seconds",
			Help:      "HTTP request latencies in seconds.",
			Labels:    labels,
		})

	// request size
	reqSizeBytes = metric.NewHistogramVec(
		&metric.HistogramVecOpts{
			Namespace: namespace,
			Name:      "http_request_size_bytes",
			Help:      "HTTP request sizes in bytes.",
			Labels:    labels,
		})

	// response size
	respSizeBytes = metric.NewHistogramVec(
		&metric.HistogramVecOpts{
			Namespace: namespace,
			Name:      "http_response_size_bytes",
			Help:      "HTTP request sizes in bytes.",
			Labels:    labels,
		})
)

// calcRequestSize returns the size of request object.
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

// Metrics returns a gin.HandlerFunc for exporting some Web metrics
func Metrics(serviceName string) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		c.Next()

		status := fmt.Sprintf("%d", c.Writer.Status())
		endpoint := c.Request.URL.Path
		method := c.Request.Method

		labels := []string{status, endpoint, method, serviceName}

		// no response content will return -1
		respSize := c.Writer.Size()
		if respSize < 0 {
			respSize = 0
		}
		curReqCount.Inc(labels...)
		defer curReqCount.Dec(labels...)
		reqCount.Inc(labels...)
		reqDuration.Observe(int64(time.Since(start).Seconds()), labels...)
		reqSizeBytes.Observe(int64(calcRequestSize(c.Request)), labels...)
		respSizeBytes.Observe(int64(respSize), labels...)
	}
}
