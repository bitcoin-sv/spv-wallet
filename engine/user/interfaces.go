package user

import (
	"context"

	paymailmodels "github.com/bitcoin-sv/spv-wallet/engine/paymail/models"
	"github.com/bitcoin-sv/spv-wallet/engine/user/usermodels"
)

// UsersRepo is a user repository
type UsersRepo interface {
	AppendAddress(ctx context.Context, userID string, newAddress *usermodels.NewAddress) error
	CreateUser(ctx context.Context, newUser *usermodels.NewUser) (*usermodels.User, error)
	AppendPaymail(ctx context.Context, userID string, newPaymail *usermodels.NewPaymail) (*paymailmodels.Paymail, error)
	GetWithPaymails(ctx context.Context, userID string) (*usermodels.User, error)
}
