package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Metrics is the metrics collector
type Metrics struct {
	gatherer   prometheus.Gatherer
	registerer prometheus.Registerer
}

// newMetrics is private to ensure that only one global-instance is created
func newMetrics() *Metrics {
	registry := prometheus.NewRegistry()
	constLabels := prometheus.Labels{"app": appName}
	registererWithLabels := prometheus.WrapRegistererWith(constLabels, registry)

	m := &Metrics{
		gatherer:   registry,
		registerer: registererWithLabels,
	}

	return m
}

// HTTPHandler will return the http.Handler for the metrics
func (m *Metrics) HTTPHandler() http.Handler {
	return promhttp.HandlerFor(metrics.gatherer, promhttp.HandlerOpts{Registry: metrics.registerer})
}
