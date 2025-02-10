package data

import (
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/rs/zerolog"
)

// APIData represents server with API endpoints
type APIData struct {
	engine engine.ClientInterface
	logger *zerolog.Logger
}

// NewAPIData creates a new server with API endpoints
func NewAPIData(engine engine.ClientInterface, log *zerolog.Logger) APIData {
	return APIData{
		engine: engine,
		logger: log,
	}
}
