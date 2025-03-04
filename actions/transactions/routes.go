package transactions

import (
	"fmt"
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

	registerARCCallback(handlersManager)
}

func registerARCCallback(handlersManager *handlers.Manager) {
	config := handlersManager.GetConfig()
	if config.ARCCallbackEnabled() {
		callbackURL, err := config.ARC.Callback.ShouldGetURL()
		if err != nil {
			panic(fmt.Sprintf(`couldn't get callback URL from configuration: %v`, err.Error()))
		}
		handlersManager.Get(handlers.GroupTransactionCallback).POST(callbackURL.Path, broadcastCallback)
	}
}
