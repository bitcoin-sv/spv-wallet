package utxos

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
	router.HTTPRouter.GET("/"+config.APIVersion+"/utxo", action.Request(router, require.Wrap(action.get)))
	router.HTTPRouter.POST("/"+config.APIVersion+"/utxo/count", action.Request(router, require.Wrap(action.count)))
	router.HTTPRouter.POST("/"+config.APIVersion+"/utxo/search", action.Request(router, require.Wrap(action.search)))
}
