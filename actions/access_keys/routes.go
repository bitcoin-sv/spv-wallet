package accesskeys

import (
	"github.com/bitcoin-sv/spv-wallet/server/handlers"
)

// RegisterRoutes creates the specific package routes
func RegisterRoutes(handlersManager *handlers.Manager) {
	old := handlersManager.Group(handlers.GroupOldAPI, "/access-key")
	old.POST("", handlers.AsUserWithXPub(oldCreate))
	old.GET("", handlers.AsUser(oldGet))
	old.DELETE("", handlers.AsUserWithXPub(oldRevoke))
	old.POST("/count", handlers.AsUser(count))
	old.POST("/search", handlers.AsUser(oldSearch))

	group := handlersManager.Group(handlers.GroupAPI, "/users/current/keys")
	group.GET("/:id", handlers.AsUser(get))
	group.POST("", handlers.AsUserWithXPub(create))
	group.DELETE("/:id", handlers.AsUserWithXPub(revoke))
	group.GET("", handlers.AsUser(search))
}
