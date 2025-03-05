package operations

import (
	v2 "github.com/bitcoin-sv/spv-wallet/engine/v2"
	"github.com/rs/zerolog"
)

// APIOperations represents server with API endpoints
type APIOperations struct {
	engine v2.Engine
	logger *zerolog.Logger
}

// NewAPIOperations creates a new server with API endpoints
func NewAPIOperations(engine v2.Engine, log *zerolog.Logger) APIOperations {
	logger := log.With().Str("api", "operations").Logger()

	return APIOperations{
		engine: engine,
		logger: &logger,
	}
}
