package accesskeys

import (
	"github.com/BuxOrg/spv-wallet/actions"
	"github.com/BuxOrg/spv-wallet/config"
	apirouter "github.com/mrz1836/go-api-router"
)

// Action is an extension of actions.Action for this package
type Action struct {
	actions.Action
}

// RegisterRoutes register all the package specific routes
func RegisterRoutes(router *apirouter.Router, appConfig *config.AppConfig, services *config.AppServices) {

	// Use the authentication middleware wrapper - this will only check for a valid xPub
	a, require := actions.NewStack(appConfig, services)
	require.Use(a.RequireAuthentication)

	// Load the actions and set the services
	action := &Action{actions.Action{AppConfig: a.AppConfig, Services: a.Services}}

	// V1 Requests
	router.HTTPRouter.GET("/"+config.APIVersion+"/access-key", action.Request(router, require.Wrap(action.get)))
	router.HTTPRouter.POST("/"+config.APIVersion+"/access-key/count", action.Request(router, require.Wrap(action.count)))
	router.HTTPRouter.GET("/"+config.APIVersion+"/access-key/search", action.Request(router, require.Wrap(action.search)))
	router.HTTPRouter.POST("/"+config.APIVersion+"/access-key/search", action.Request(router, require.Wrap(action.search)))
	router.HTTPRouter.POST("/"+config.APIVersion+"/access-key", action.Request(router, require.Wrap(action.create)))
	router.HTTPRouter.DELETE("/"+config.APIVersion+"/access-key", action.Request(router, require.Wrap(action.revoke)))

}
