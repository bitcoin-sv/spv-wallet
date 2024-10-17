package outlines

import "context"

// Service is a service for creating transaction outlines.
type Service interface {
	Create(ctx context.Context, spec *TransactionSpec) (*Transaction, error)
}
