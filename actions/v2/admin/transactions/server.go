package transactions

import (
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/rs/zerolog"
)

// APIAdminTransactions represents server with admin API endpoints
type APIAdminTransactions struct {
	engine engine.ClientInterface
	logger *zerolog.Logger
}

// NewAPIAdminTransactions creates a new APIAdminTransactions
func NewAPIAdminTransactions(engine engine.ClientInterface, logger *zerolog.Logger) APIAdminTransactions {
	return APIAdminTransactions{
		engine: engine,
		logger: logger,
	}
}
