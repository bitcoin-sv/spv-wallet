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

	old.POST("", handlers.AsUserWithXPub(newTransaction))
	old.POST("/record", handlers.AsUserWithXPub(record))

	group := handlersManager.Group(handlers.GroupAPI, "/transactions")
	group.GET(":id", handlers.AsUser(getByID))
	group.PATCH(":id", handlers.AsUser(updateTransactionMetadata))
	group.GET("", handlers.AsUser(transactions))

	group.POST("/drafts", handlers.AsUserWithXPub(newTransactionDraft))
	group.POST("", handlers.AsUserWithXPub(recordTransaction))

	handlersManager.Get(handlers.GroupTransactionCallback).POST(config.BroadcastCallbackRoute, broadcastCallback)
}
