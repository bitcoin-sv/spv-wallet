package chainmodels

// QueryTXOutcome represents the outcome of a transaction query.
type QueryTXOutcome int

// QueryTXOutcome values sorted from most "negative" to most "positive".
const (
	QueryTxOutcomeFailed QueryTXOutcome = iota
	QueryTXOutcomeRejected
	QueryTXOutcomeNotFound
	QueryTXOutcomeSuccess
)
