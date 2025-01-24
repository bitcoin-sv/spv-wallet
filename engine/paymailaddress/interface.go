package paymailaddress

import (
	"context"

	"github.com/bitcoin-sv/spv-wallet/engine/database"
)

// Service is a component that provides methods for working with paymail address.
type Service interface {
	HasPaymailAddress(ctx context.Context, userID string, address string) (bool, error)
	GetDefaultPaymailAddress(ctx context.Context, userID string) (string, error)
}

type PaymailRepo interface {
	// FindForUser returns a paymail by alias and domain for given user.
	FindForUser(ctx context.Context, alias, domain, userID string) (*database.Paymail, error)
	// GetDefault returns a default paymail for user.
	GetDefault(ctx context.Context, userID string) (*database.Paymail, error)
}
