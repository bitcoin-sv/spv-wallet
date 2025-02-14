package admin

import (
	"github.com/bitcoin-sv/spv-wallet/actions/v2/admin/users"
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/rs/zerolog"
)

// APIAdmin represents server with API endpoints
type APIAdmin struct {
	users.APIAdminUsers
}

// NewAPIAdmin creates a new APIAdmin
func NewAPIAdmin(spvWalletEngine engine.ClientInterface, logger *zerolog.Logger) APIAdmin {
	return APIAdmin{
		users.NewAPIAdminUsers(spvWalletEngine, logger),
	}
}
