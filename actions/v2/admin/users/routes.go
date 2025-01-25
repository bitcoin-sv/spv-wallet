package users

import (
	"github.com/bitcoin-sv/spv-wallet/server/handlers"
	routes "github.com/bitcoin-sv/spv-wallet/server/handlers"
)

// RegisterRoutes creates the specific package routes
func RegisterRoutes(handlersManager *routes.Manager) {
	group := handlersManager.Group(routes.GroupAPIV2, "admin/users")
	group.POST("", handlers.AsAdmin(create))
	group.GET(":id", handlers.AsAdmin(get))

	group.POST(":id/paymails", handlers.AsAdmin(addPaymail))
}
