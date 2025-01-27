package v2

import (
	"github.com/bitcoin-sv/spv-wallet/actions/v2/admin"
	"github.com/bitcoin-sv/spv-wallet/actions/v2/operations"
	"github.com/bitcoin-sv/spv-wallet/actions/v2/users"
	"github.com/bitcoin-sv/spv-wallet/server/handlers"
)

// Register collects all the action's routes and registers them using the handlersManager
func Register(handlersManager *handlers.Manager) {
	users.RegisterRoutes(handlersManager)
	operations.RegisterRoutes(handlersManager)

	admin.Register(handlersManager)
}
