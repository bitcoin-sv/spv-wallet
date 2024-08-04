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

// OldAccessKeysHandler creates the specific package routes
func OldAccessKeysHandler(appConfig *config.AppConfig, services *config.AppServices) routes.OldAPIEndpointsFunc {
	action := &Action{actions.Action{AppConfig: appConfig, Services: services}}

	oldAPIEndpoints := routes.OldAPIEndpointsFunc(func(router *gin.RouterGroup) {
		accessKeyGroup := router.Group("/access-key")
		accessKeyGroup.POST("", action.oldCreate)
		accessKeyGroup.GET("", action.oldGet)
		accessKeyGroup.DELETE("", action.oldRevoke)
		accessKeyGroup.POST("/count", action.count)
		accessKeyGroup.GET("/search", action.oldSearch)
		accessKeyGroup.POST("/search", action.oldSearch)
	})

	return oldAPIEndpoints
}

// NewHandler creates the specific package routes
func NewHandler(appConfig *config.AppConfig, services *config.AppServices) routes.APIEndpointsFunc {
	action := &Action{actions.Action{AppConfig: appConfig, Services: services}}

	apiEndpoints := routes.APIEndpointsFunc(func(router *gin.RouterGroup) {
		accessKeyGroup := router.Group("/users/current/keys")
		accessKeyGroup.GET("/:id", action.get)
		accessKeyGroup.POST("", action.create)
		accessKeyGroup.DELETE("/:id", action.revoke)
		accessKeyGroup.GET("", action.searchTest)
	})

	return apiEndpoints
}
