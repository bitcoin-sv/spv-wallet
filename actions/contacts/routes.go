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
		contactGroup := router.Group("/contact")
		contactGroup.POST("", action.create)
		contactGroup.PATCH("", action.update)
	})

	return apiEndpoints
}
