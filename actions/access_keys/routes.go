package accesskeys

import (
	"github.com/bitcoin-sv/spv-wallet/server/handlers"
)

// RegisterRoutes creates the specific package routes
func RegisterRoutes(handlersManager *handlers.Manager) {
	group := handlersManager.Group(handlers.GroupAPI, "/users/current/keys")
	group.GET("/:id", handlers.AsUser(get))
	group.POST("", handlers.AsUser(create))
	group.DELETE("/:id", handlers.AsUser(revoke))
	group.GET("", handlers.AsUser(search))
}
