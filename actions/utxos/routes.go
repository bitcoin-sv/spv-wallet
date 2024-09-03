package utxos

import (
	"github.com/bitcoin-sv/spv-wallet/server/handlers"
	routes "github.com/bitcoin-sv/spv-wallet/server/handlers"
)

// RegisterRoutes creates the specific package routes
func RegisterRoutes(handlersManager *routes.Manager) {
	group := handlersManager.Group(routes.GroupOldAPI, "/utxo")
	group.GET("", handlers.AsUser(get))
	group.POST("/count", handlers.AsUser(count))
	group.POST("/search", handlers.AsUser(search))
}
