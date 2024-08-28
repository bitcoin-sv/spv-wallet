package chainstate

// broadcastQuestionableErrors are a list of errors that are not good broadcast responses,
// but need to be checked differently
var broadcastQuestionableErrors = []string{
	"missing inputs", // The transaction has been sent to at least 1 Bitcoin node but parent transaction was not found. This status means that inputs are currently missing, but the transaction is not yet rejected.
}

/*
	TXN_ALREADY_KNOWN (suppressed - returns as success: true)
	TXN_ALREADY_IN_MEMPOOL (suppressed - returns as success: true)
	TXN_MEMPOOL_CONFLICT
	NON_FINAL_POOL_FULL
	TOO_LONG_NON_FINAL_CHAIN
	BAD_TXNS_INPUTS_TOO_LARGE
	BAD_TXNS_INPUTS_SPENT
	NON_BIP68_FINAL
	TOO_LONG_VALIDATION_TIME
	BAD_TXNS_NONSTANDARD_INPUTS
	ABSURDLY_HIGH_FEE
	DUST
	TX_FEE_TOO_LOW
*/
