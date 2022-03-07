package accessKeys

import (
	"github.com/BuxOrg/bux-server/actions"
	"github.com/BuxOrg/bux-server/config"
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
	router.HTTPRouter.DELETE("/"+config.CurrentMajorVersion+"/access_key", router.Request(require.Wrap(action.revoke)))
	router.HTTPRouter.GET("/"+config.CurrentMajorVersion+"/access_key", router.Request(require.Wrap(action.get)))
	router.HTTPRouter.GET("/"+config.CurrentMajorVersion+"/access_keys", router.Request(require.Wrap(action.search)))
	router.HTTPRouter.POST("/"+config.CurrentMajorVersion+"/access_key", router.Request(require.Wrap(action.create)))
	router.HTTPRouter.POST("/"+config.CurrentMajorVersion+"/access_keys", router.Request(require.Wrap(action.search)))
}
