package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func NewMetrics(prefix string) Metrics {
	registry := prometheus.NewRegistry()
	registerer := prometheus.WrapRegistererWithPrefix(prefix, registry)

	registerer.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
	)

	return &prometheusWrapper{
		registry:   registry,
		registerer: registerer,
	}
}

type prometheusWrapper struct {
	registry   *prometheus.Registry
	registerer prometheus.Registerer
}

func (m *prometheusWrapper) NewCounter(name, description string) Counter {
	collector := prometheus.NewCounter(prometheus.CounterOpts{
		Name: name,
		Help: description,
	})
	m.registerer.MustRegister(collector)
	return collector
}

func (m *prometheusWrapper) NewGauge(name, description string) Gauge {
	collector := prometheus.NewGauge(prometheus.GaugeOpts{
		Name: name,
		Help: description,
	})
	m.registerer.MustRegister(collector)
	return collector
}

func (m *prometheusWrapper) NewHistogram(name, description string) Histogram {
	collector := prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: name,
		Help: description,
	})
	m.registerer.MustRegister(collector)
	return collector
}

func (m *prometheusWrapper) HttpHandler() http.Handler {
	return promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{Registry: m.registerer})
}
