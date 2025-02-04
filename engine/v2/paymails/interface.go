package paymails

import (
	"context"

	"github.com/bitcoin-sv/spv-wallet/engine/v2/paymails/paymailsmodels"
)

// PaymailRepo is a paymail repository
type PaymailRepo interface {
	Create(ctx context.Context, newPaymail *paymailsmodels.NewPaymail) (*paymailsmodels.Paymail, error)
	Find(ctx context.Context, alias, domain string) (*paymailsmodels.Paymail, error)
	// FindForUser returns a paymail by alias and domain for given user.
	FindForUser(ctx context.Context, alias, domain, userID string) (*paymailsmodels.Paymail, error)
	// GetDefault returns a default paymail for user.
	GetDefault(ctx context.Context, userID string) (*paymailsmodels.Paymail, error)
}

// UsersService is a user domain service
type UsersService interface {
	Exists(ctx context.Context, userID string) (bool, error)
}

// DomainChecker is an interface for checking domain.
type DomainChecker interface {
	CheckDomain(domain string) error
}
