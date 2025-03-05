package data

import (
	v2 "github.com/bitcoin-sv/spv-wallet/engine/v2"
	"github.com/rs/zerolog"
)

// APIData represents server with API endpoints
type APIData struct {
	engine v2.Engine
	logger *zerolog.Logger
}

// NewAPIData creates a new server with API endpoints
func NewAPIData(engine v2.Engine, log *zerolog.Logger) APIData {
	logger := log.With().Str("api", "data").Logger()

	return APIData{
		engine: engine,
		logger: &logger,
	}
}
