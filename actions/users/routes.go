package users

import (
	"github.com/bitcoin-sv/spv-wallet/server/handlers"
	routes "github.com/bitcoin-sv/spv-wallet/server/handlers"
)

// RegisterRoutes creates the specific package routes in RESTful style
func RegisterRoutes(handlersManager *routes.Manager) {
	group := handlersManager.Group(routes.GroupAPI, "/users/current")
	group.GET("", handlers.AsUser(get))
	group.PATCH("", handlers.AsUser(update))
}
