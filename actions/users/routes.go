package users

import (
	"github.com/bitcoin-sv/spv-wallet/server/handlers"
	routes "github.com/bitcoin-sv/spv-wallet/server/handlers"
)

// RegisterRoutes creates the specific package routes in RESTful style
func RegisterRoutes(handlersManager *routes.Manager) {
	old := handlersManager.Group(routes.GroupOldAPI, "/xpub")
	old.GET("", handlers.AsUser(oldGet))
	old.PATCH("", handlers.AsUser(oldUpdate))

	group := handlersManager.Group(routes.GroupAPI, "/users/current")
	group.GET("", handlers.AsUser(get))
	group.PATCH("", handlers.AsUser(update))

	v2 := handlersManager.Group(routes.GroupAPIV2, "/users")
	v2.GET("balance", handlers.AsUser(balance))
	v2.GET("operations", handlers.AsUser(operations))
}
