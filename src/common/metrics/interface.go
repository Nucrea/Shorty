package metrics

type Counter interface {
	Inc()
}

type Gauge interface {
	Set(float64)
	// Inc()
	// Dec()
}

type Histogram interface {
	Observe(float64)
	ObserveWithLabel(value float64, label string)
}

type Meter interface {
	NewCounter(name, description string) Counter
	NewGauge(name, description string) Gauge
	NewHistogram(name, description string, buckets []float64) Histogram
}
