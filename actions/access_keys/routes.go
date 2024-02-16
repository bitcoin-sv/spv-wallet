package accesskeys

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
func NewHandler(appConfig *config.AppConfig, services *config.AppServices) routes.ApiEndpointsFunc {
	action := &Action{actions.Action{AppConfig: appConfig, Services: services}}

	apiEndpoints := routes.ApiEndpointsFunc(func(router *gin.RouterGroup) {
		accessKeyGroup := router.Group("/access-key")
		accessKeyGroup.POST("", action.create)
		accessKeyGroup.GET("", action.get)
		accessKeyGroup.DELETE("", action.revoke)
		accessKeyGroup.POST("/count", action.count)
		accessKeyGroup.GET("/search", action.search)
		accessKeyGroup.POST("/search", action.search)
	})

	return apiEndpoints
}
