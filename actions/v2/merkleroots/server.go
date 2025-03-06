package merkleroots

import (
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/rs/zerolog"
)

// APIMerkleRoots represents server with API endpoints
type APIMerkleRoots struct {
	engine engine.ClientInterface
	logger *zerolog.Logger
}

// NewAPIMerkleRoots creates a new server with API endpoints
func NewAPIMerkleRoots(engine engine.ClientInterface, log *zerolog.Logger) APIMerkleRoots {
	logger := log.With().Str("api", "merkleroots").Logger()

	return APIMerkleRoots{
		engine: engine,
		logger: &logger,
	}
}
