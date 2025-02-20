package operations

import (
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/rs/zerolog"
)

// APIOperations represents server with API endpoints
type APIOperations struct {
	engine engine.ClientInterface
	logger *zerolog.Logger
}

// NewAPIOperations creates a new server with API endpoints
func NewAPIOperations(engine engine.ClientInterface, log *zerolog.Logger) APIOperations {
	logger := log.With().Str("api", "operations").Logger()

	return APIOperations{
		engine: engine,
		logger: &logger,
	}
}
