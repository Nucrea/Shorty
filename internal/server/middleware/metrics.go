package middleware

import (
	"shorty/internal/common/metrics"
	"time"

	"github.com/gin-gonic/gin"
)

func Metrics(meter metrics.Meter) gin.HandlerFunc {
	errorsMetric := meter.NewCounter("http_requests_errors", "HTTP Errors (status >= 500)")
	totalMetric := meter.NewCounter("http_requests_total", "HTTP Requests")
	latencyMetric := meter.NewGauge("http_requests_latency", "HTTP Requests latency")

	return func(ctx *gin.Context) {
		start := time.Now()
		defer func() {
			duration := time.Now().Sub(start)
			latencyMetric.Set(float64(duration.Milliseconds()))

			totalMetric.Inc()
			if ctx.Writer.Status() >= 500 {
				errorsMetric.Inc()
			}
		}()

		ctx.Next()
	}
}
