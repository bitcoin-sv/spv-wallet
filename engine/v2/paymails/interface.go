package paymails

import (
	"context"

	"github.com/bitcoin-sv/spv-wallet/engine/v2/paymails/paymailsmodels"
)

// PaymailRepo is a paymail repository
type PaymailRepo interface {
	Create(ctx context.Context, newPaymail *paymailsmodels.NewPaymail) (*paymailsmodels.Paymail, error)
	Get(ctx context.Context, alias, domain string) (*paymailsmodels.Paymail, error)
}

// UsersService is a user domain service
type UsersService interface {
	Exists(ctx context.Context, userID string) (bool, error)
}
