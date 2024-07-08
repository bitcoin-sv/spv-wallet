package sharedconfig

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
func NewHandler(appConfig *config.AppConfig, services *config.AppServices) (routes.OldAPIEndpointsFunc, routes.APIEndpointsFunc) {
	action := &Action{actions.Action{AppConfig: appConfig, Services: services}}

	oldSharedConfigEndpoints := routes.OldAPIEndpointsFunc(func(router *gin.RouterGroup) {
		group := router.Group("/shared-config")
		group.GET("", action.oldGet)
	})

	sharedConfigEndpoints := routes.APIEndpointsFunc(func(router *gin.RouterGroup) {
		group := router.Group("/configs/shared")
		group.GET("", action.get)
	})

	return oldSharedConfigEndpoints, sharedConfigEndpoints
}
