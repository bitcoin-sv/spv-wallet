package transactions

import (
	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/bitcoin-sv/spv-wallet/server/handlers"
)

// RegisterRoutes creates the specific package routes
func RegisterRoutes(handlersManager *handlers.Manager) {
	group := handlersManager.Group(handlers.GroupAPI, "/transactions")
	group.GET(":id", handlers.AsUser(getByID))
	group.PATCH(":id", handlers.AsUser(updateTransactionMetadata))
	group.GET("", handlers.AsUser(transactions))

	group.POST("/drafts", handlers.AsUser(newTransactionDraft))
	group.POST("", handlers.AsUser(recordTransaction))

	handlersManager.Get(handlers.GroupTransactionCallback).POST(config.BroadcastCallbackRoute, broadcastCallback)

	if handlersManager.GetFeatureFlags().NewTransactionFlowEnabled {
		v2 := handlersManager.Group(handlers.GroupAPIV2, "/transactions")
		v2.POST("/outlines", handlers.AsUser(transactionOutlines))
		v2.POST("", handlers.AsUser(transactionRecordOutline))
	}
}
