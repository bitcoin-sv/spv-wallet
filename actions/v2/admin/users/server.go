package users

import (
	v2 "github.com/bitcoin-sv/spv-wallet/engine/v2"
	"github.com/rs/zerolog"
)

// APIAdminUsers represents server with admin API endpoints
type APIAdminUsers struct {
	engine v2.Engine
	logger *zerolog.Logger
}

// NewAPIAdminUsers creates a new APIAdminUsers
func NewAPIAdminUsers(engine v2.Engine, logger *zerolog.Logger) APIAdminUsers {
	return APIAdminUsers{
		engine: engine,
		logger: logger,
	}
}
