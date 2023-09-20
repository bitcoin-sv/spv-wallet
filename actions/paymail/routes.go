package pmail

import (
	"github.com/BuxOrg/bux-server/actions"
	"github.com/BuxOrg/bux-server/config"
	//apirouter "github.com/mrz1836/go-api-router"
	"github.com/mrz1836/go-logger"
	"github.com/newrelic/go-agent/v3/integrations/nrhttprouter"
)

// Action is an extension of actions.Action for this package
type Action struct {
	actions.Action
}

// RegisterRoutes register all the package specific routes
// func RegisterRoutes(router *nrhttprouter.Router, appConfig *config.AppConfig, services *config.AppServices) {
func RegisterRoutes(router *nrhttprouter.Router, appConfig *config.AppConfig, services *config.AppServices) {
	// Use the authentication middleware wrapper
	a, requireAdmin := actions.NewStack(appConfig, services)
	requireAdmin.Use(a.RequireAdminAuthentication)

	// Register the custom Paymail routes
	services.Bux.GetPaymailConfig().RegisterRoutes(router.HTTPRouter)

	// Create the action
	action := &Action{actions.Action{AppConfig: a.AppConfig, Services: a.Services}}

	// V1 Requests
	router.POST("/"+config.CurrentMajorVersion+"/paymail", requireAdmin.Wrap(action.create))
	router.DELETE("/"+config.CurrentMajorVersion+"/paymail", requireAdmin.Wrap(action.delete))

	if appConfig.Debug {
		logger.Data(2, logger.DEBUG, "registered paymail routes and model")
	}
}
