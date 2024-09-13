package draft

import "context"

// Service is a service for creating draft transactions.
type Service interface {
	Create(ctx context.Context, spec *TransactionSpec) (*Transaction, error)
}
