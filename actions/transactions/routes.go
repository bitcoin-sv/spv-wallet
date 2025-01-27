package transactions

import (
	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/bitcoin-sv/spv-wallet/server/handlers"
)

// RegisterRoutes creates the specific package routes
func RegisterRoutes(handlersManager *handlers.Manager) {
	old := handlersManager.Group(handlers.GroupOldAPI, "/transaction")
	old.GET("", handlers.AsUser(get))
	old.PATCH("", handlers.AsUser(update))
	old.POST("/count", handlers.AsUser(count))
	old.GET("/search", handlers.AsUser(search))
	old.POST("/search", handlers.AsUser(search))

	old.POST("", handlers.AsUser(newTransaction))
	old.POST("/record", handlers.AsUser(record))

	group := handlersManager.Group(handlers.GroupAPI, "/transactions")
	group.GET(":id", handlers.AsUser(getByID))
	group.PATCH(":id", handlers.AsUser(updateTransactionMetadata))
	group.GET("", handlers.AsUser(transactions))

	group.POST("/drafts", handlers.AsUser(newTransactionDraft))
	group.POST("", handlers.AsUser(recordTransaction))

	handlersManager.Get(handlers.GroupTransactionCallback).POST(config.BroadcastCallbackRoute, broadcastCallback)

	if handlersManager.GetFeatureFlags().V2 {
		v2 := handlersManager.Group(handlers.GroupAPIV2, "/transactions")
		v2.POST("/outlines", handlers.AsUser(transactionOutlines))
		v2.POST("", handlers.AsUser(transactionRecordOutline))
	}
}
