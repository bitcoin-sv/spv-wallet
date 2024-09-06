package accesskeys

import (
	"github.com/bitcoin-sv/spv-wallet/server/handlers"
)

// RegisterRoutes creates the specific package routes
func RegisterRoutes(handlersManager *handlers.Manager) {
	old := handlersManager.Group(handlers.GroupOldAPI, "/access-key")
	old.POST("", handlers.AsUser(oldCreate))
	old.GET("", handlers.AsUser(oldGet))
	old.DELETE("", handlers.AsUser(oldRevoke))
	old.POST("/count", handlers.AsUser(count))
	old.POST("/search", handlers.AsUser(oldSearch))

	group := handlersManager.Group(handlers.GroupAPI, "/users/current/keys")
	group.GET("/:id", handlers.AsUser(get))
	group.POST("", handlers.AsUser(create))
	group.DELETE("/:id", handlers.AsUser(revoke))
	group.GET("", handlers.AsUser(search))
}
