package users

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

// OldUsersHandler creates the specific package routes
func OldUsersHandler(appConfig *config.AppConfig, services *config.AppServices) routes.OldAPIEndpointsFunc {
	action := &Action{actions.Action{AppConfig: appConfig, Services: services}}

	oldAPIEndpoints := routes.OldAPIEndpointsFunc(func(router *gin.RouterGroup) {
		xpubGroup := router.Group("/xpub")
		xpubGroup.GET("", action.oldGet)
		xpubGroup.PATCH("", action.oldUpdate)
	})

	return oldAPIEndpoints
}

// NewHandler creates the specific package routes in RESTful style
func NewHandler(appConfig *config.AppConfig, services *config.AppServices) routes.APIEndpointsFunc {
	action := &Action{actions.Action{AppConfig: appConfig, Services: services}}

	apiEndpoints := routes.APIEndpointsFunc(func(router *gin.RouterGroup) {
		xpubGroup := router.Group("/users/current")
		xpubGroup.GET("", action.get)
		xpubGroup.PATCH("", action.update)
	})

	return apiEndpoints
}
