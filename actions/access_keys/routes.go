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
func NewHandler(appConfig *config.AppConfig, services *config.AppServices) routes.OldAPIEndpointsFunc {
	action := &Action{actions.Action{AppConfig: appConfig, Services: services}}

	apiEndpoints := routes.OldAPIEndpointsFunc(func(router *gin.RouterGroup) {
		accessKeyGroup := router.Group("/access-key")
		accessKeyGroup.POST("", action.create)
		accessKeyGroup.GET("", action.get)
		accessKeyGroup.DELETE("", action.revoke)
		accessKeyGroup.POST("/count", action.count)
		accessKeyGroup.GET("/search", action.search)
		accessKeyGroup.POST("/search", action.search)

		accessKeyGroup2 := router.Group("/users/current/keys")
		accessKeyGroup2.GET("/:id", action.get2)
		accessKeyGroup2.POST("", action.create2)
		accessKeyGroup2.DELETE("/:id", action.revoke2)
		// TODO: accessKeyGroup2.GET("", action.search)
	})

	return apiEndpoints
}
