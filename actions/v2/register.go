package v2

import (
	"github.com/bitcoin-sv/spv-wallet/actions/v2/admin"
	"github.com/bitcoin-sv/spv-wallet/actions/v2/transactions"
	"github.com/bitcoin-sv/spv-wallet/server/handlers"
)

// Register collects all the action's routes and registers them using the handlersManager
func Register(handlersManager *handlers.Manager) {
	transactions.RegisterRoutes(handlersManager)

	admin.Register(handlersManager)
}
