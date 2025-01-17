package actions

import (
	accesskeys "github.com/bitcoin-sv/spv-wallet/actions/access_keys"
	"github.com/bitcoin-sv/spv-wallet/actions/admin"
	"github.com/bitcoin-sv/spv-wallet/actions/base"
	"github.com/bitcoin-sv/spv-wallet/actions/contacts"
	"github.com/bitcoin-sv/spv-wallet/actions/merkleroots"
	"github.com/bitcoin-sv/spv-wallet/actions/paymails"
	"github.com/bitcoin-sv/spv-wallet/actions/sharedconfig"
	"github.com/bitcoin-sv/spv-wallet/actions/transactions"
	"github.com/bitcoin-sv/spv-wallet/actions/users"
	"github.com/bitcoin-sv/spv-wallet/actions/utxos"
	"github.com/bitcoin-sv/spv-wallet/server/handlers"
)

// Register collects all the action's routes and registers them using the handlersManager
func Register(handlersManager *handlers.Manager) {
	admin.RegisterRoutes(handlersManager)
	base.RegisterRoutes(handlersManager)
	accesskeys.RegisterRoutes(handlersManager)
	transactions.RegisterRoutes(handlersManager)
	utxos.RegisterRoutes(handlersManager)
	users.RegisterRoutes(handlersManager)
	paymails.RegisterRoutes(handlersManager)
	sharedconfig.RegisterRoutes(handlersManager)
	merkleroots.RegisterRoutes(handlersManager)
	contacts.RegisterRoutes(handlersManager)
}
