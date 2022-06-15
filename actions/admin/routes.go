package admin

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

	// Use the authentication middleware wrapper - this will only check for a valid admin
	a, require := actions.NewStack(appConfig, services)
	require.Use(a.RequireAdminAuthentication)

	// Load the actions and set the services
	action := &Action{actions.Action{AppConfig: a.AppConfig, Services: a.Services}}

	// V1 Requests
	router.HTTPRouter.GET("/"+config.CurrentMajorVersion+"/admin/stats", action.Request(router, require.Wrap(action.stats)))
	router.HTTPRouter.GET("/"+config.CurrentMajorVersion+"/admin/status", action.Request(router, require.Wrap(action.status)))
	router.HTTPRouter.POST("/"+config.CurrentMajorVersion+"/admin/access-keys/search", action.Request(router, require.Wrap(action.accessKeysSearch)))
	router.HTTPRouter.POST("/"+config.CurrentMajorVersion+"/admin/access-keys/count", action.Request(router, require.Wrap(action.accessKeysCount)))
	router.HTTPRouter.POST("/"+config.CurrentMajorVersion+"/admin/block-headers/search", action.Request(router, require.Wrap(action.blockHeadersSearch)))
	router.HTTPRouter.POST("/"+config.CurrentMajorVersion+"/admin/block-headers/count", action.Request(router, require.Wrap(action.blockHeadersCount)))
	router.HTTPRouter.POST("/"+config.CurrentMajorVersion+"/admin/destinations/search", action.Request(router, require.Wrap(action.destinationsSearch)))
	router.HTTPRouter.POST("/"+config.CurrentMajorVersion+"/admin/destinations/count", action.Request(router, require.Wrap(action.destinationsCount)))
	router.HTTPRouter.POST("/"+config.CurrentMajorVersion+"/admin/paymails/search", action.Request(router, require.Wrap(action.paymailAddressesSearch)))
	router.HTTPRouter.POST("/"+config.CurrentMajorVersion+"/admin/paymails/count", action.Request(router, require.Wrap(action.paymailAddressesCount)))
	router.HTTPRouter.POST("/"+config.CurrentMajorVersion+"/admin/transactions/search", action.Request(router, require.Wrap(action.transactionsSearch)))
	router.HTTPRouter.POST("/"+config.CurrentMajorVersion+"/admin/transactions/count", action.Request(router, require.Wrap(action.transactionsCount)))
	router.HTTPRouter.POST("/"+config.CurrentMajorVersion+"/admin/transactions/record", action.Request(router, require.Wrap(action.transactionRecord)))
	router.HTTPRouter.POST("/"+config.CurrentMajorVersion+"/admin/utxos/search", action.Request(router, require.Wrap(action.utxosSearch)))
	router.HTTPRouter.POST("/"+config.CurrentMajorVersion+"/admin/utxos/count", action.Request(router, require.Wrap(action.utxosCount)))
	router.HTTPRouter.POST("/"+config.CurrentMajorVersion+"/admin/xpubs/search", action.Request(router, require.Wrap(action.xpubsSearch)))
	router.HTTPRouter.POST("/"+config.CurrentMajorVersion+"/admin/xpubs/count", action.Request(router, require.Wrap(action.xpubsCount)))
}
