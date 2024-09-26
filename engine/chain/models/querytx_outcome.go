package chainmodels

type QueryTXOutcome int

const (
	QueryTxOutcomeFailed QueryTXOutcome = iota
	QueryTXOutcomeRejected
	QueryTXOutcomeNotFound
	QueryTXOutcomeSuccess
)
