package metrics

import (
	"fmt"
	"time"

	enginemetrics "github.com/bitcoin-sv/spv-wallet/engine/metrics"
	"github.com/prometheus/client_golang/prometheus"
)

// RequestMetrics is the metrics for the http requests
type RequestMetrics struct {
	requestsTotal   *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
}

func registerRequestMetrics(collector enginemetrics.Collector) *RequestMetrics {
	requestsTotal := collector.RegisterCounterVec(requestCounterName, "method", "path", "status", "classification")
	requestDuration := collector.RegisterHistogramVec(requestDurationSecName, "method", "path")

	return &RequestMetrics{
		requestsTotal:   requestsTotal,
		requestDuration: requestDuration,
	}
}

// Track will return a RequestTracker to track the request
func (m *RequestMetrics) Track(method, path string) *RequestTracker {
	return &RequestTracker{
		method:  method,
		path:    path,
		metrics: m,
	}
}

// RequestTracker is used to track the duration and status of a request
type RequestTracker struct {
	method    string
	path      string
	startTime time.Time
	metrics   *RequestMetrics
}

// Start will start the tracking of the request
func (r *RequestTracker) Start() {
	r.startTime = time.Now()
}

// End will end the tracking of the request
func (r *RequestTracker) End(status int) {
	r.writeCounter(status, r.path)
	r.writeDuration()
}

// EndWithNoRoute will end the tracking of the request with a 404 status
func (r *RequestTracker) EndWithNoRoute() {
	// This is a safeguard against attacks where the server is flooded with requests having unique paths,
	// which would lead to the creation of a large number of metrics
	r.writeCounter(404, "UNKNOWN_ROUTE")
}

func (r *RequestTracker) writeCounter(status int, path string) {
	r.metrics.requestsTotal.WithLabelValues(r.method, path, fmt.Sprint(status), requestClassification(status)).Inc()
}

func (r *RequestTracker) writeDuration() {
	r.metrics.requestDuration.WithLabelValues(r.method, r.path).Observe(time.Since(r.startTime).Seconds())
}

func requestClassification(status int) string {
	if status >= 200 && status < 400 {
		return "success"
	}
	return "failure"
}
