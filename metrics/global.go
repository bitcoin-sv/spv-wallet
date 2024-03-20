package metrics

import (
	enginemetrics "github.com/bitcoin-sv/spv-wallet/engine/metrics"
	"github.com/gin-gonic/gin"
)

var metrics *Metrics

// EnableMetrics will enable the metrics for the application
func EnableMetrics() enginemetrics.Collector {
	var collector enginemetrics.Collector
	metrics, collector = newMetrics()
	return collector
}

// Get will return the metrics if enabled
func Get() (m *Metrics, enabled bool) {
	return metrics, metrics != nil
}

// SetupGin will register the metrics with the gin engine
// NOTE: Remember to add the metrics.NoRoute function to ginEngine.NoRoute
func SetupGin(ginEngine *gin.Engine) {
	if metrics != nil {
		ginEngine.Use(requestMetricsMiddleware())
		ginEngine.GET("/metrics", gin.WrapH(metrics.HTTPHandler()))
	}
}
