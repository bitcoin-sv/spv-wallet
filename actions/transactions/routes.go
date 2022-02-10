package transactions

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

	// Load the actions and set the services
	action := &Action{actions.Action{AppConfig: a.AppConfig, Services: a.Services}}

	// V1 Requests
	router.HTTPRouter.POST("/"+config.CurrentMajorVersion+"/transactions/record", router.Request(require.Wrap(action.record)))
	router.HTTPRouter.POST("/"+config.CurrentMajorVersion+"/transactions/new", router.Request(require.Wrap(action.newTransaction)))
	router.HTTPRouter.GET("/"+config.CurrentMajorVersion+"/transactions", router.Request(require.Wrap(action.list)))
	router.HTTPRouter.GET("/"+config.CurrentMajorVersion+"/transaction", router.Request(require.Wrap(action.get)))
}
