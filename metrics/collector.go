package metrics

import (
	enginemetrics "github.com/bitcoin-sv/spv-wallet/engine/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

// PrometheusCollector is a collector for Prometheus metrics. It should implement spvwalletmodels.Collector.
type PrometheusCollector struct {
	reg prometheus.Registerer
}

// NewPrometheusCollector creates a new PrometheusCollector.
func NewPrometheusCollector(reg prometheus.Registerer) enginemetrics.Collector {
	return &PrometheusCollector{reg: reg}
}

// RegisterGauge creates a new Gauge and registers it with the collector.
func (c *PrometheusCollector) RegisterGauge(name string) prometheus.Gauge {
	g := prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: name,
			Help: "Gauge of " + name,
		},
	)
	c.reg.MustRegister(g)
	return g
}

// RegisterGaugeVec creates a new GaugeVec and registers it with the collector.
func (c *PrometheusCollector) RegisterGaugeVec(name string, labels ...string) *prometheus.GaugeVec {
	g := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: name,
			Help: "GaugeVec of " + name,
		},
		labels,
	)
	c.reg.MustRegister(g)
	return g
}

// RegisterHistogramVec creates a new HistogramVec and registers it with the collector.
func (c *PrometheusCollector) RegisterHistogramVec(name string, labels ...string) *prometheus.HistogramVec {
	h := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: name,
			Help: "HistogramVec of " + name,
		},
		labels,
	)
	c.reg.MustRegister(h)
	return h
}
