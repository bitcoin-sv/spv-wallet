package utxos

import (
	"github.com/bitcoin-sv/spv-wallet/server/handlers"
	routes "github.com/bitcoin-sv/spv-wallet/server/handlers"
)

// NewHandler creates the specific package routes
func NewHandler(handlersManager *routes.Manager) {
	roGroup := handlersManager.Group(routes.GroupOldAPI, "/utxo")
	roGroup.GET("", handlers.AsUser(get))
	roGroup.POST("/count", handlers.AsUser(count))
	roGroup.POST("/search", handlers.AsUser(search))
}
