package contacts

import (
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/rs/zerolog"
)

// APIAdminContacts represents server with admin API endpoints
type APIAdminContacts struct {
	engine engine.ClientInterface
	logger *zerolog.Logger
}

// NewAPIAdminContacts creates a new APIAdminUsers
func NewAPIAdminContacts(engine engine.ClientInterface, logger *zerolog.Logger) APIAdminContacts {
	return APIAdminContacts{
		engine: engine,
		logger: logger,
	}
}
