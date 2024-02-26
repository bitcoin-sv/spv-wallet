package destinations

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
func NewHandler(appConfig *config.AppConfig, services *config.AppServices) (routes.BasicEndpointsFunc, routes.APIEndpointsFunc) {
	action := &Action{actions.Action{AppConfig: appConfig, Services: services}}

	basicEndpoints := routes.BasicEndpointsFunc(func(router *gin.RouterGroup) {
		basicDestinationGroup := router.Group("/destination")
		basicDestinationGroup.GET("", action.get)
		basicDestinationGroup.POST("/count", action.count)
		basicDestinationGroup.GET("/search", action.search)
		basicDestinationGroup.POST("/search", action.search)
	})

	apiEndpoints := routes.APIEndpointsFunc(func(router *gin.RouterGroup) {
		apiDestinationGroup := router.Group("/destination")
		apiDestinationGroup.POST("", action.create)
		apiDestinationGroup.PATCH("", action.update)
	})

	return basicEndpoints, apiEndpoints
}
