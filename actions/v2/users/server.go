package users

import (
	"context"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/rs/zerolog"
)

type usersService interface {
	Remove(ctx context.Context, userID string) error
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
