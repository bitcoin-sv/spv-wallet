package merkleroots

import (
	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/bitcoin-sv/spv-wallet/server/handlers"
)

func RegisterRoutes(appConfig *config.AppConfig, handlersManager *handlers.Manager) {
	group := handlersManager.Group(handlers.GroupAPI, "/merkleroots")
	group.GET("", handlers.AsUserWithAppConfig(get, appConfig))
}
