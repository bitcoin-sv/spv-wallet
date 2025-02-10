package transactions

import (
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/rs/zerolog"
)

// APITransactions represents server with API endpoints
type APITransactions struct {
	engine engine.ClientInterface
	logger *zerolog.Logger
}

// NewAPITransactions creates a new server with API endpoints
func NewAPITransactions(engine engine.ClientInterface, log *zerolog.Logger) APITransactions {
	return APITransactions{
		engine: engine,
		logger: log,
	}
}
