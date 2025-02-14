package v2

import (
	"github.com/bitcoin-sv/spv-wallet/actions/v2/swagger"
	"github.com/bitcoin-sv/spv-wallet/server/handlers"
)

// RegisterNonOpenAPIRoutes collects all the action's routes that aren't part of the Open API documentation and registers them using the handlersManager.
func RegisterNonOpenAPIRoutes(handlersManager *handlers.Manager) {
	swagger.RegisterRoutes(handlersManager)
}
