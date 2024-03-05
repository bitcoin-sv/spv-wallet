package metrics

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	enginemetrics "github.com/bitcoin-sv/spv-wallet/engine/metrics"
)

// Metrics is the metrics collector
type Metrics struct {
	gatherer     prometheus.Gatherer
	registerer   prometheus.Registerer
	httpRequests *RequestMetrics
}

// newMetrics is private to ensure that only one global-instance is created
func newMetrics() (*Metrics, enginemetrics.Collector) {
	registry := prometheus.NewRegistry()
	constLabels := prometheus.Labels{"app": appName}
	registererWithLabels := prometheus.WrapRegistererWith(constLabels, registry)

	collector := newPrometheusCollector(registererWithLabels)

	m := &Metrics{
		gatherer:     registry,
		registerer:   registererWithLabels,
		httpRequests: registerRequestMetrics(collector),
	}

	return m, collector
}

// HTTPHandler will return the http.Handler for the metrics
func (m *Metrics) HTTPHandler() http.Handler {
	return promhttp.HandlerFor(metrics.gatherer, promhttp.HandlerOpts{Registry: metrics.registerer})
}
