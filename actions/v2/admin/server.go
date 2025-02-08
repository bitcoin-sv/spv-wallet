package admin

import "github.com/bitcoin-sv/spv-wallet/actions/v2/admin/users"

// APIAdmin represents server with API endpoints
type APIAdmin struct {
	users.APIAdminUsers
}
