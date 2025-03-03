package transactions

import (
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

	cfg := handlersManager.GetConfig()
	if cfg.ARC.Callback.Enabled {
		callbackPath := cfg.ARC.Callback.MustGetURL().Path
		handlersManager.Get(handlers.GroupTransactionCallback).POST(callbackPath, broadcastCallback)
	}
}
