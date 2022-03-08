package pmail

import (
	"context"

	"github.com/BuxOrg/bux"
	"github.com/BuxOrg/bux-server/actions"
	"github.com/BuxOrg/bux-server/config"
	apirouter "github.com/mrz1836/go-api-router"
	"github.com/mrz1836/go-logger"
)

// Action is an extension of actions.Action for this package
type Action struct {
	actions.Action
}

// RegisterRoutes register all the package specific routes
func RegisterRoutes(router *apirouter.Router, appConfig *config.AppConfig, services *config.AppServices) {

	// Use the authentication middleware wrapper
	a, requireAdmin := actions.NewStack(appConfig, services)
	requireAdmin.Use(a.RequireAdminAuthentication)

	// Register the routes
	services.Bux.PaymailServerConfig().RegisterRoutes(router)

	// Add the additional models
	// todo: ideally, this should be in Services or on-load (cyclical dep issue)
	if err := services.Bux.AddModels(context.Background(), true, &bux.PaymailAddress{}); err != nil {
		// todo: handle this error (avoid using a panic)
		panic(err)
	}

	action := &Action{actions.Action{AppConfig: a.AppConfig, Services: a.Services}}

	// V1 Requests
	router.HTTPRouter.POST("/"+config.CurrentMajorVersion+"/paymail", router.Request(requireAdmin.Wrap(action.create)))
	router.HTTPRouter.DELETE("/"+config.CurrentMajorVersion+"/paymail", router.Request(requireAdmin.Wrap(action.delete)))

	if appConfig.Debug {
		logger.Data(2, logger.DEBUG, "registered paymail routes and model")
	}
}
