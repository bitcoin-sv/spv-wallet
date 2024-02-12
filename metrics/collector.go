package metrics

import (
	buxmetrics "github.com/BuxOrg/bux/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

// PrometheusCollector is a collector for Prometheus metrics. It should implement buxmetrics.Collector.
type PrometheusCollector struct {
	reg prometheus.Registerer
}

// NewPrometheusCollector creates a new PrometheusCollector.
func NewPrometheusCollector(reg prometheus.Registerer) buxmetrics.Collector {
	return &PrometheusCollector{reg: reg}
}

// RegisterGauge creates a new Gauge and registers it with the collector.
func (c *PrometheusCollector) RegisterGauge(name string) buxmetrics.GaugeInterface {
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
func (c *PrometheusCollector) RegisterGaugeVec(name string, labels ...string) buxmetrics.GaugeVecInterface {
	g := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: name,
			Help: "GaugeVec of " + name,
		},
		labels,
	)
	c.reg.MustRegister(g)
	return &GaugeVecWrapper{g}
}

// RegisterHistogramVec creates a new HistogramVec and registers it with the collector.
func (c *PrometheusCollector) RegisterHistogramVec(name string, labels ...string) buxmetrics.HistogramVecInterface {
	h := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: name,
			Help: "HistogramVec of " + name,
		},
		labels,
	)
	c.reg.MustRegister(h)
	return &HistogramVecWrapper{h}
}

// GaugeVecWrapper is a wrapper for prometheus.GaugeVec
type GaugeVecWrapper struct {
	*prometheus.GaugeVec
}

// WithLabelValues returns a Gauge with the given label values
func (g *GaugeVecWrapper) WithLabelValues(lvs ...string) buxmetrics.GaugeInterface {
	return g.GaugeVec.WithLabelValues(lvs...)
}

// HistogramVecWrapper is a wrapper for prometheus.HistogramVec
type HistogramVecWrapper struct {
	*prometheus.HistogramVec
}

// WithLabelValues returns a Histogram with the given label values
func (h *HistogramVecWrapper) WithLabelValues(lvs ...string) buxmetrics.HistogramInterface {
	return h.HistogramVec.WithLabelValues(lvs...)
}
