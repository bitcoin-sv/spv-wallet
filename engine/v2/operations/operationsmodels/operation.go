package operationsmodels

import "time"

// Operation represents a user's operation on with underlying transaction.
type Operation struct {
	TxID   string
	UserID string

	CreatedAt time.Time

	Counterparty string
	Type         string
	Value        int64

	TxStatus string
}
