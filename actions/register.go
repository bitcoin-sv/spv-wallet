package actions

import (
	accesskeys "github.com/bitcoin-sv/spv-wallet/actions/access_keys"
	"github.com/bitcoin-sv/spv-wallet/actions/admin"
	"github.com/bitcoin-sv/spv-wallet/actions/base"
	"github.com/bitcoin-sv/spv-wallet/actions/contacts"
	"github.com/bitcoin-sv/spv-wallet/actions/destinations"
	"github.com/bitcoin-sv/spv-wallet/actions/sharedconfig"
	"github.com/bitcoin-sv/spv-wallet/actions/transactions"
	"github.com/bitcoin-sv/spv-wallet/actions/users"
	"github.com/bitcoin-sv/spv-wallet/actions/utxos"
	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/bitcoin-sv/spv-wallet/server/handlers"
)

// Register collects all the action's routes and registers them using the handlersManager
func Register(appConfig *config.AppConfig, handlersManager *handlers.Manager) {
	admin.RegisterRoutes(handlersManager)
	base.RegisterRoutes(handlersManager)
	accesskeys.RegisterRoutes(handlersManager)
	destinations.RegisterRoutes(handlersManager)
	transactions.RegisterRoutes(handlersManager)
	utxos.RegisterRoutes(handlersManager)
	users.RegisterRoutes(handlersManager)
	sharedconfig.RegisterRoutes(handlersManager)
	if appConfig.ExperimentalFeatures.PikeContactsEnabled {
		contacts.RegisterRoutes(handlersManager)
	}
}
