package outlines

import (
	"context"

	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction"
)

// PaymailAddressService is a component that provides methods for working with paymail address.
type PaymailAddressService interface {
	HasPaymailAddress(ctx context.Context, userID string, address string) (bool, error)
	GetDefaultPaymailAddress(ctx context.Context, userID string) (string, error)
}

// Service is a service for creating transaction outlines.
type Service interface {
	Create(ctx context.Context, spec *TransactionSpec) (*Transaction, error)
}

// Transaction represents a transaction outline.
type Transaction struct {
	BEEF        string
	Annotations transaction.Annotations
}
