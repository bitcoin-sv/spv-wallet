package contacts

import (
	"github.com/bitcoin-sv/spv-wallet/actions"
	"github.com/bitcoin-sv/spv-wallet/config"
	apirouter "github.com/mrz1836/go-api-router"
)

// Action is an extension of actions.Action for this package
type Action struct {
	actions.Action
}

// RegisterRoutes register all the package specific routes
func RegisterRoutes(router *apirouter.Router, appConfig *config.AppConfig, services *config.AppServices) {
	a, require := actions.NewStack(appConfig, services)
	require.Use(a.RequireAuthentication)

	// Use the authentication middleware wrapper - this will only check for a valid xPub
	aBasic, requireBasic := actions.NewStack(appConfig, services)
	requireBasic.Use(aBasic.RequireBasicAuthentication)

	// Load the actions and set the services
	action := &Action{actions.Action{AppConfig: a.AppConfig, Services: a.Services}}

	router.HTTPRouter.POST("/"+config.APIVersion+"/contact", action.Request(router, requireBasic.Wrap(action.create)))
}
