package sharedconfig

import (
	"github.com/bitcoin-sv/spv-wallet/server/handlers"
)

// RegisterRoutes creates the specific package routes
func RegisterRoutes(handlersManager *handlers.Manager) {
	old := handlersManager.Group(handlers.GroupOldAPI, "/shared-config")
	old.GET("", handlers.AsAdminOrUser(oldGet))

	group := handlersManager.Group(handlers.GroupAPI, "/configs/shared")
	group.GET("", handlers.AsAdminOrUser(get))
}
