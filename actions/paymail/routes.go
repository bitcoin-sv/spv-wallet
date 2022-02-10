package pmail

import (
	"context"

	"github.com/BuxOrg/bux-server/config"
	apirouter "github.com/mrz1836/go-api-router"
	"github.com/mrz1836/go-logger"
)

// RegisterRoutes register all the package specific routes
func RegisterRoutes(router *apirouter.Router, appConfig *config.AppConfig, services *config.AppServices) {

	// Use the authentication middleware wrapper
	// a, _ := actions.NewStack(appConfig, services)
	// require.Use(a.RequireAuthentication)

	// Register the routes
	services.Bux.PaymailServerConfig().RegisterRoutes(router)

	// Add the additional models
	// todo: ideally, this should be in Services or on-load (cyclical dep issue)
	if err := services.Bux.AddModels(context.Background(), true, &PaymailAddress{}); err != nil {
		// todo: handle this error (avoid using a panic)
		panic(err)
	}

	if appConfig.Debug {
		logger.Data(2, logger.DEBUG, "registered paymail routes and model")
	}
}
