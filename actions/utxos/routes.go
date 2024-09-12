package utxos

import (
	"github.com/bitcoin-sv/spv-wallet/server/handlers"
	routes "github.com/bitcoin-sv/spv-wallet/server/handlers"
)

// RegisterRoutes creates the specific package routes
func RegisterRoutes(handlersManager *routes.Manager) {
	old := handlersManager.Group(routes.GroupOldAPI, "/utxo")
	old.GET("", handlers.AsUser(get))
	old.POST("/count", handlers.AsUser(count))
	old.POST("/search", handlers.AsUser(oldSearch))

	group := handlersManager.Group(routes.GroupAPI, "/utxos")
	group.GET("", handlers.AsUser(search))
}
