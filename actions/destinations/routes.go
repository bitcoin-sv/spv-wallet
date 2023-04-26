package destinations

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
	// Use the authentication middleware wrapper
	a, require := actions.NewStack(appConfig, services)
	require.Use(a.RequireAuthentication)

	// Use the authentication middleware wrapper - this will only check for a valid xPub
	aBasic, requireBasic := actions.NewStack(appConfig, services)
	requireBasic.Use(aBasic.RequireBasicAuthentication)

	// Load the actions and set the services
	action := &Action{actions.Action{AppConfig: a.AppConfig, Services: a.Services}}

	// V1 Requests
	router.HTTPRouter.GET("/"+config.CurrentMajorVersion+"/destination", action.Request(router, requireBasic.Wrap(action.get)))
	router.HTTPRouter.POST("/"+config.CurrentMajorVersion+"/destination/count", action.Request(router, requireBasic.Wrap(action.count)))
	router.HTTPRouter.GET("/"+config.CurrentMajorVersion+"/destination/search", action.Request(router, requireBasic.Wrap(action.search)))
	router.HTTPRouter.POST("/"+config.CurrentMajorVersion+"/destination", action.Request(router, require.Wrap(action.create)))
	router.HTTPRouter.PATCH("/"+config.CurrentMajorVersion+"/destination", action.Request(router, require.Wrap(action.update)))
}
