package users

import (
	"context"

	primitives "github.com/bitcoin-sv/go-sdk/primitives/ec"
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/users/usersmodels"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/rs/zerolog"
)

type usersService interface {
	Create(ctx context.Context, newUser *usersmodels.NewUser) (*usersmodels.User, error)
	Remove(ctx context.Context, userID string) error
	Exists(ctx context.Context, userID string) (bool, error)
	GetByID(ctx context.Context, userID string) (*usersmodels.User, error)
	GetIDByPubKey(ctx context.Context, pubKey string) (string, error)
	GetPubKey(ctx context.Context, userID string) (*primitives.PublicKey, error)
	GetBalance(ctx context.Context, userID string) (bsv.Satoshis, error)
}

// APIUsers represents server with API endpoints
type APIUsers struct {
	usersService usersService
	logger       *zerolog.Logger
}

// NewAPIUsers creates a new server with API endpoints
func NewAPIUsers(engine engine.ClientInterface, log *zerolog.Logger) APIUsers {
	logger := log.With().Str("api", "users").Logger()

	return APIUsers{
		usersService: engine.UsersService(),
		logger:       &logger,
	}
}
