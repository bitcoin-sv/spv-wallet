package paymails

import (
	"github.com/bitcoin-sv/spv-wallet/server/handlers"
	routes "github.com/bitcoin-sv/spv-wallet/server/handlers"
)

// RegisterRoutes creates the specific package routes in RESTful style
func RegisterRoutes(handlersManager *routes.Manager) {
	group := handlersManager.Group(routes.GroupAPI, "/paymails")
	group.GET("", handlers.AsUser(paymailAddressesSearch))
}
