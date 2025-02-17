package users

import (
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/rs/zerolog"
)

// APIUsers represents server with API endpoints
type APIUsers struct {
	engine engine.ClientInterface
	logger *zerolog.Logger
}

// NewAPIUsers creates a new server with API endpoints
func NewAPIUsers(engine engine.ClientInterface, log *zerolog.Logger) APIUsers {
	logger := log.With().Str("api", "users").Logger()

	return APIUsers{
		engine: engine,
		logger: &logger,
	}
}
