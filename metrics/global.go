package metrics

import (
	spvmetrics "github.com/bitcoin-sv/spv-wallet/engine/metrics"
)

var metrics *Metrics

// EnableMetrics will enable the metrics for the application
func EnableMetrics() spvmetrics.Collector {
	metrics = newMetrics()
	return NewPrometheusCollector(metrics.registerer)
}

// Get will return the metrics if enabled
func Get() (m *Metrics, enabled bool) {
	return metrics, metrics != nil
}
