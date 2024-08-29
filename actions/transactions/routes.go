package transactions

import (
	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/bitcoin-sv/spv-wallet/server/handlers"
	"github.com/bitcoin-sv/spv-wallet/server/middleware"
)

// NewHandler creates the specific package routes
func NewHandler(handlersManager *handlers.Manager) {
	old := handlersManager.Group(handlers.GroupOldAPI, "/transaction")
	old.GET("", handlers.AsUser(get))
	old.PATCH("", handlers.AsUser(update))
	old.POST("/count", handlers.AsUser(count))
	old.GET("/search", handlers.AsUser(search))
	old.POST("/search", handlers.AsUser(search))

	old.POST("", middleware.RequireSignature, handlers.AsUser(newTransaction))
	old.POST("/record", middleware.RequireSignature, handlers.AsUser(record))

	group := handlersManager.Group(handlers.GroupAPI, "/transactions")
	group.GET(":id", handlers.AsUser(getByID))
	group.PATCH(":id", handlers.AsUser(updateTransactionMetadata))
	group.GET("", handlers.AsUser(transactions))

	group.POST("/drafts", middleware.RequireSignature, handlers.AsUser(newTransactionDraft))
	group.POST("", middleware.RequireSignature, handlers.AsUser(recordTransaction))

	handlersManager.Get(handlers.GroupTransactionCallback).POST(config.BroadcastCallbackRoute, broadcastCallback)
}
