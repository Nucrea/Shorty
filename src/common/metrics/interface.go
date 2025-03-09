package metrics

import "net/http"

type Counter interface {
	Inc()
}

type Gauge interface {
	Set(float64)
	Inc()
	Dec()
}

type Histogram interface {
	Observe(float64)
}

type Metrics interface {
	NewCounter(name, description string) Counter
	NewGauge(name, description string) Gauge
	NewHistogram(name, description string) Histogram
	HttpHandler() http.Handler
}
