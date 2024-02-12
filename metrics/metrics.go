package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Metrics struct {
	gatherer   prometheus.Gatherer
	registerer prometheus.Registerer
}

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

func (m *Metrics) HttpHandler() http.Handler {
	return promhttp.HandlerFor(metrics.gatherer, promhttp.HandlerOpts{Registry: metrics.registerer})
}
