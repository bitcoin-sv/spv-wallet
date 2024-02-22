package metrics

import "github.com/prometheus/client_golang/prometheus"

// Collector is an interface that is used to register metrics
type Collector interface {
	RegisterGauge(name string) prometheus.Gauge
	RegisterGaugeVec(name string, labels ...string) *prometheus.GaugeVec
	RegisterHistogramVec(name string, labels ...string) *prometheus.HistogramVec
}
