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

// OldContactsHandler creates the specific package routes
func OldContactsHandler(appConfig *config.AppConfig, services *config.AppServices) routes.OldAPIEndpointsFunc {
	action := &Action{actions.Action{AppConfig: appConfig, Services: services}}

	apiEndpoints := routes.OldAPIEndpointsFunc(func(router *gin.RouterGroup) {
		group := router.Group("/contact")
		group.PUT("/:paymail", action.oldUpsert)

		group.PATCH("/accepted/:paymail", action.oldAccept)
		group.PATCH("/rejected/:paymail", action.oldReject)
		group.PATCH("/confirmed/:paymail", action.oldConfirm)
		group.PATCH("/unconfirmed/:paymail", action.oldUnconfirm)

		group.POST("search", action.search)
	})

	return apiEndpoints
}

// NewHandler creates the specific package routes
func NewHandler(appConfig *config.AppConfig, services *config.AppServices) (routes.APIEndpointsFunc, routes.APIEndpointsFunc) {
	action := &Action{actions.Action{AppConfig: appConfig, Services: services}}

	contactsAPIEndpoints := routes.APIEndpointsFunc(func(router *gin.RouterGroup) {
		group := router.Group("/contacts")
		group.PUT("/:paymail", action.upsertContact)
		group.POST("/:paymail/confirmation", action.confirmContact)
		group.DELETE("/:paymail/confirmation", action.unconfirmContact)

		group.GET("", action.getContacts)
		group.GET(":paymail", action.getContactByPaymail)
	})

	invitationsAPIEndpoints := routes.APIEndpointsFunc(func(router *gin.RouterGroup) {
		group := router.Group("/invitations")

		group.POST("/:paymail", action.acceptInvitations)
		group.DELETE("/:paymail", action.rejectInvitation)

	})
	return contactsAPIEndpoints, invitationsAPIEndpoints
}
