package merkleroots

import (
	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/bitcoin-sv/spv-wallet/server/handlers"
)

// RegisterRoutes creates the specific package routes
func RegisterRoutes(appConfig *config.AppConfig, handlersManager *handlers.Manager) {
	group := handlersManager.Group(handlers.GroupAPI, "/merkleroots")
	group.GET("", handlers.AsUserWithAppConfig(get, appConfig))
}
