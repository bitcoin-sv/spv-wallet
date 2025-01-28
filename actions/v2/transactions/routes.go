package transactions

import "github.com/bitcoin-sv/spv-wallet/server/handlers"

// RegisterRoutes creates the specific package routes
func RegisterRoutes(handlersManager *handlers.Manager) {
	if handlersManager.GetFeatureFlags().V2 {
		group := handlersManager.Group(handlers.GroupAPIV2, "/transactions")
		group.POST("/outlines", handlers.AsUser(transactionOutlines))
		group.POST("", handlers.AsUser(transactionRecordOutline))
	}
}
