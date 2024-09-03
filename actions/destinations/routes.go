package destinations

import (
	"github.com/bitcoin-sv/spv-wallet/server/handlers"
	"github.com/bitcoin-sv/spv-wallet/server/middleware"
)

// RegisterRoutes creates the specific package routes
func RegisterRoutes(handlersManager *handlers.Manager) {
	group := handlersManager.Group(handlers.GroupOldAPI, "/destination")
	group.GET("", handlers.AsUser(get))
	group.POST("/count", handlers.AsUser(count))
	group.GET("/search", handlers.AsUser(search))
	group.POST("/search", handlers.AsUser(search))

	group.POST("", middleware.RequireSignature, handlers.AsUser(create))
	group.PATCH("", middleware.RequireSignature, handlers.AsUser(update))
}
