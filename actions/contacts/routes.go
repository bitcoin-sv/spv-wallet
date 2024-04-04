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

		group.POST("search", action.search)
	})

	return apiEndpoints
}
