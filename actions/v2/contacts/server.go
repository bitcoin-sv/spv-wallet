package contacts

import (
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/rs/zerolog"
)

// APIContacts represents server with API endpoints
type APIContacts struct {
	engine engine.ClientInterface
	logger *zerolog.Logger
}

// NewAPIContacts creates a new server with API endpoints
func NewAPIContacts(engine engine.ClientInterface, log *zerolog.Logger) APIContacts {
	logger := log.With().Str("api", "contacts").Logger()

	return APIContacts{
		engine: engine,
		logger: &logger,
	}
}
