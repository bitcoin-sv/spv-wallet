package xpubs

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
		xpubGroup := router.Group("/xpub")
		xpubGroup.GET("", action.get)
		xpubGroup.PATCH("", action.update)

		xpubGroup2 := router.Group("/users/current")
		xpubGroup2.GET("", action.get2)
		xpubGroup2.PATCH("", action.update2)
	})

	return apiEndpoints
}
