package users

import (
	"context"

	"github.com/bitcoin-sv/spv-wallet/engine/v2/users/usersmodels"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
)

// UserRepo is an interface for users repository.
type UserRepo interface {
	Exists(ctx context.Context, userID string) (bool, error)
	GetIDByPubKey(ctx context.Context, pubKey string) (string, error)
	Get(ctx context.Context, userID string) (*usersmodels.User, error)
	Create(ctx context.Context, newUser *usersmodels.NewUser) (*usersmodels.User, error)
	GetBalance(ctx context.Context, userID string, bucket string) (bsv.Satoshis, error)
}

//// AddressRepo is an interface for addresses repository.
//type AddressRepo interface {
//	Create(ctx context.Context, newAddress *domainmodels.NewAddress) error
//}
//
//// PaymailRepo is an interface for paymails repository.
//type PaymailRepo interface {
//	Create(ctx context.Context, newPaymail *domainmodels.NewPaymail) (*domainmodels.Paymail, error)
//}
