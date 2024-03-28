package metrics

const domainPrefix = "bsv_"

const (
	verifyMerkleRootsHistogramName = domainPrefix + "verify_merkle_roots_histogram"
	recordTransactionHistogramName = domainPrefix + "record_transaction_histogram"
	queryTransactionHistogramName  = domainPrefix + "query_transaction_histogram"
	addContactHistogramName        = domainPrefix + "add_contact_histogram"
)

const (
	cronHistogramName          = domainPrefix + "cron_histogram"
	cronLastExecutionGaugeName = domainPrefix + "cron_last_execution_gauge"
)

const (
	statsGaugeName = domainPrefix + "stats_total"
)
