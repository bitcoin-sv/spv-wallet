package operations

import (
	"github.com/bitcoin-sv/spv-wallet/server/handlers"
	routes "github.com/bitcoin-sv/spv-wallet/server/handlers"
)

// RegisterRoutes creates the specific package routes
func RegisterRoutes(handlersManager *routes.Manager) {
	group := handlersManager.Group(routes.GroupAPIV2, "/operations")
	group.GET("search", handlers.AsUser(search))
}
