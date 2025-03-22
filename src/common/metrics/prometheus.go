package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var _ Meter = (*prometheusMetrics)(nil)

func NewPrometheus(prefix string) *prometheusMetrics {
	registry := prometheus.NewRegistry()
	registerer := prometheus.WrapRegistererWithPrefix(prefix, registry)

	registerer.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	return &prometheusMetrics{
		registry:   registry,
		registerer: registerer,
	}
}

type prometheusMetrics struct {
	registry   *prometheus.Registry
	registerer prometheus.Registerer
}

func (m *prometheusMetrics) NewCounter(name, description string) Counter {
	collector := prometheus.NewCounter(prometheus.CounterOpts{
		Name: name,
		Help: description,
	})
	m.registerer.MustRegister(collector)
	return collector
}

func (m *prometheusMetrics) NewGauge(name, description string) Gauge {
	collector := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: name,
		Help: description,
	})
	m.registerer.MustRegister(collector)
	return collector
}

func (m *prometheusMetrics) NewHistogram(name, description string, buckets []float64) Histogram {
	collector := prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    name,
		Help:    description,
		Buckets: buckets,
	})
	m.registerer.MustRegister(collector)
	return collector
}

func (m *prometheusMetrics) HttpHandler() http.Handler {
	return promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{Registry: m.registerer})
}
