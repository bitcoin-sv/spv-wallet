package admin

import (
	"github.com/bitcoin-sv/spv-wallet/actions/v2/admin/users"
	v2 "github.com/bitcoin-sv/spv-wallet/engine/v2"
	"github.com/rs/zerolog"
)

// APIAdmin represents server with API endpoints
type APIAdmin struct {
	users.APIAdminUsers
}

// NewAPIAdmin creates a new APIAdmin
func NewAPIAdmin(spvWalletEngine v2.Engine, logger *zerolog.Logger) APIAdmin {
	return APIAdmin{
		users.NewAPIAdminUsers(spvWalletEngine.UsersService(), spvWalletEngine.PaymailsService(), logger),
	}
}
