package transactions

import "github.com/bitcoin-sv/spv-wallet/server/handlers"

// RegisterRoutes creates the specific package routes
func RegisterRoutes(handlersManager *handlers.Manager) {
	if handlersManager.GetFeatureFlags().NewTransactionFlowEnabled {
		v2 := handlersManager.Group(handlers.GroupAPIV2, "/transactions")
		v2.POST("/outlines", handlers.AsUser(transactionOutlines))
		v2.POST("", handlers.AsUser(transactionRecordOutline))
	}
}
