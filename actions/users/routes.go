package users

import (
	"github.com/bitcoin-sv/spv-wallet/server/handlers"
	routes "github.com/bitcoin-sv/spv-wallet/server/handlers"
	"github.com/bitcoin-sv/spv-wallet/server/middleware"
)

// NewHandler creates the specific package routes in RESTful style
func NewHandler(handlersManager *routes.Manager) {
	old := handlersManager.Group(routes.GroupOldAPI, "/xpub")
	old.GET("", handlers.AsUser(oldGet))
	old.PATCH("", middleware.RequireSignature, handlers.AsUser(oldUpdate))

	group := handlersManager.Group(routes.GroupAPI, "/users/current")
	group.GET("", handlers.AsUser(get))
	group.PATCH("", middleware.RequireSignature, handlers.AsUser(update))
}
