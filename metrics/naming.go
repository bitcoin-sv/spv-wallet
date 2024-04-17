package metrics

const appName = "spv-wallet"

const (
	requestMetricBaseName  = "http_request"
	requestCounterName     = requestMetricBaseName + "_total"
	requestDurationSecName = requestMetricBaseName + "_duration_seconds"
)
