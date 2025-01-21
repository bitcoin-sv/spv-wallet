package users

import (
	"context"

	"github.com/bitcoin-sv/spv-wallet/engine/domainmodels"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
)

// UserRepo is an interface for users repository.
type UserRepo interface {
	Exists(ctx context.Context, userID string) (bool, error)
	GetByPubKey(ctx context.Context, pubKey string) (*domainmodels.User, error)
	Get(ctx context.Context, userID string) (*domainmodels.User, error)
	Create(ctx context.Context, newUser *domainmodels.NewUser) (*domainmodels.User, error)
	GetBalance(ctx context.Context, userID string, bucket string) (bsv.Satoshis, error)
}

// AddressRepo is an interface for addresses repository.
type AddressRepo interface {
	Create(ctx context.Context, newAddress *domainmodels.NewAddress) error
}

// PaymailRepo is an interface for paymails repository.
type PaymailRepo interface {
	Create(ctx context.Context, newPaymail *domainmodels.NewPaymail) (*domainmodels.Paymail, error)
}
