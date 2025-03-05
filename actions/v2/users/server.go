package users

import (
	v2 "github.com/bitcoin-sv/spv-wallet/engine/v2"
	"github.com/rs/zerolog"
)

// APIUsers represents server with API endpoints
type APIUsers struct {
	engine v2.Engine
	logger *zerolog.Logger
}

// NewAPIUsers creates a new server with API endpoints
func NewAPIUsers(engine v2.Engine, log *zerolog.Logger) APIUsers {
	logger := log.With().Str("api", "users").Logger()

	return APIUsers{
		engine: engine,
		logger: &logger,
	}
}
