package admin

import (
	"github.com/bitcoin-sv/spv-wallet/actions/v2/admin/transactions"
	"github.com/bitcoin-sv/spv-wallet/actions/v2/admin/users"
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/rs/zerolog"
)

// APIAdmin represents server with API endpoints
type APIAdmin struct {
	users.APIAdminUsers
	transactions.APIAdminTransactions
}

// NewAPIAdmin creates a new APIAdmin
func NewAPIAdmin(spvWalletEngine engine.ClientInterface, logger *zerolog.Logger) APIAdmin {
	return APIAdmin{
		users.NewAPIAdminUsers(spvWalletEngine, logger),
		transactions.NewAPIAdminTransactions(spvWalletEngine, logger),
	}
}
