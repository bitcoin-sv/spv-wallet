package transactions

import (
	"github.com/bitcoin-sv/spv-wallet/server/handlers"
	routes "github.com/bitcoin-sv/spv-wallet/server/handlers"
)

// RegisterRoutes creates the specific package routes
func RegisterRoutes(handlersManager *routes.Manager) {
	group := handlersManager.Group(handlers.GroupAPIV2, "/transactions")
	group.POST("/outlines", handlers.AsUser(transactionOutlines))
	group.POST("", handlers.AsUser(recordOutline))
}
