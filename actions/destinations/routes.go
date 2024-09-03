package destinations

import (
	"github.com/bitcoin-sv/spv-wallet/server/handlers"
)

// RegisterRoutes creates the specific package routes
func RegisterRoutes(handlersManager *handlers.Manager) {
	group := handlersManager.Group(handlers.GroupOldAPI, "/destination")
	group.GET("", handlers.AsUser(get))
	group.POST("/count", handlers.AsUser(count))
	group.GET("/search", handlers.AsUser(search))
	group.POST("/search", handlers.AsUser(search))

	group.POST("", handlers.AsUser(create))
	group.PATCH("", handlers.AsUser(update))
}
