package users

import (
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/rs/zerolog"
)

// APIAdminUsers represents server with admin API endpoints
type APIAdminUsers struct {
	engine engine.ClientInterface
	logger *zerolog.Logger
}

// NewAPIAdminUsers creates a new APIAdminUsers
func NewAPIAdminUsers(engine engine.ClientInterface, logger *zerolog.Logger) APIAdminUsers {
	return APIAdminUsers{
		engine: engine,
		logger: logger,
	}
}
