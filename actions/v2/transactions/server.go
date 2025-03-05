package transactions

import (
	v2 "github.com/bitcoin-sv/spv-wallet/engine/v2"
	"github.com/rs/zerolog"
)

// APITransactions represents server with API endpoints
type APITransactions struct {
	engine v2.Engine
	logger *zerolog.Logger
}

// NewAPITransactions creates a new server with API endpoints
func NewAPITransactions(engine v2.Engine, log *zerolog.Logger) APITransactions {
	logger := log.With().Str("api", "transactions").Logger()

	return APITransactions{
		engine: engine,
		logger: &logger,
	}
}
