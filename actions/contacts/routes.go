package contacts

import (
	"github.com/bitcoin-sv/spv-wallet/actions"
	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/bitcoin-sv/spv-wallet/server/routes"
	"github.com/gin-gonic/gin"
)

// Action is an extension of actions.Action for this package
type Action struct {
	actions.Action
}

// NewHandler creates the specific package routes
func NewHandler(appConfig *config.AppConfig, services *config.AppServices) routes.APIEndpointsFunc {
	action := &Action{actions.Action{AppConfig: appConfig, Services: services}}

	apiEndpoints := routes.APIEndpointsFunc(func(router *gin.RouterGroup) {
		group := router.Group("/contact")
		group.PUT("/:paymail", action.upsert)

		group.PATCH("/accepted/:paymail", action.accept)
		group.PATCH("/rejected/:paymail", action.reject)
		group.PATCH("/confirmed/:paymail", action.confirm)
		group.PATCH("/unconfirmed/:paymail", action.unconfirm)

		group.POST("search", action.search)
	})

	return apiEndpoints
}

// NewContactsHandler creates the specific package routes
func NewContactsHandler(appConfig *config.AppConfig, services *config.AppServices) (routes.APIEndpointsFunc, routes.APIEndpointsFunc) {
	action := &Action{actions.Action{AppConfig: appConfig, Services: services}}

	contactsAPIEndpoints := routes.APIEndpointsFunc(func(router *gin.RouterGroup) {
		group := router.Group("/contacts")
		group.PUT("/:paymail", action.upsertContact)
		group.PATCH("/:paymail/confirmation", action.confirmContact)
		group.PATCH("/:paymail/non-confirmation", action.unconfirmContact)

		group.GET("", action.getContacts)
		group.GET(":paymail", action.getContactsByID)
	})

	invitationsAPIEndpoints := routes.APIEndpointsFunc(func(router *gin.RouterGroup) {
		group := router.Group("/invitations")

		group.POST("/:paymail", action.acceptInvitations)
		group.DELETE("/:paymail", action.rejectInvitation)

	})
	return contactsAPIEndpoints, invitationsAPIEndpoints
}
