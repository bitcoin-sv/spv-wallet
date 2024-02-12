package metrics

import (
	buxmetrics "github.com/BuxOrg/bux/metrics"
)

var metrics *Metrics

// EnableMetrics will enable the metrics for the application
func EnableMetrics() buxmetrics.Collector {
	metrics = newMetrics()
	return NewPrometheusCollector(metrics.registerer)
}

// Get will return the metrics if enabled
func Get() (m *Metrics, enabled bool) {
	return metrics, metrics != nil
}
