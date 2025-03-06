package users

import (
	"context"

	"github.com/bitcoin-sv/spv-wallet/engine/v2/paymails/paymailsmodels"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/users/usersmodels"
	"github.com/rs/zerolog"
)

type UsersService interface {
	Create(ctx context.Context, newUser *usersmodels.NewUser) (*usersmodels.User, error)
	GetByID(ctx context.Context, userID string) (*usersmodels.User, error)
}

type PaymailsService interface {
	Create(ctx context.Context, newPaymail *paymailsmodels.NewPaymail) (*paymailsmodels.Paymail, error)
}

// APIAdminUsers represents server with admin API endpoints
type APIAdminUsers struct {
	users    UsersService
	paymails PaymailsService
	logger   *zerolog.Logger
}

// NewAPIAdminUsers creates a new APIAdminUsers
func NewAPIAdminUsers(users UsersService, paymails PaymailsService, logger *zerolog.Logger) APIAdminUsers {
	if logger == nil {
		panic("nil logger implementation provided as argument")
	}
	if paymails == nil {
		panic("nil paymail service implementation provided as argument")
	}
	if users == nil {
		panic("nil user service implementation provided as argument")
	}

	return APIAdminUsers{
		paymails: paymails,
		users:    users,
		logger:   logger,
	}
}
