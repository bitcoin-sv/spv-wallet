/*
Package main is the core service layer for the BUX Server
*/
package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/BuxOrg/bux-server/config"
	"github.com/BuxOrg/bux-server/dictionary"
	_ "github.com/BuxOrg/bux-server/docs"
	"github.com/BuxOrg/bux-server/logging"
	"github.com/BuxOrg/bux-server/server"
)

// main method starts everything for the BUX Server
// @title           BUX: Server
// @version         v0.10.0
// @securityDefinitions.apikey bux-auth-xpub
// @in header
// @name bux-auth-xpub
func main() {
	defaultLogger := logging.GetDefaultLogger()

	// Load the Application Configuration
	appConfig, err := config.Load(defaultLogger)
	if err != nil {
		defaultLogger.Fatal().Msgf(dictionary.GetInternalMessage(dictionary.ErrorLoadingConfig), err.Error())
		return
	}

	// Validate configuration (before services have been loaded)
	if err = appConfig.Validate(); err != nil {
		defaultLogger.Fatal().Msgf(dictionary.GetInternalMessage(dictionary.ErrorLoadingConfig), err.Error())
		return
	}

	// Load the Application Services
	var services *config.AppServices
	if services, err = appConfig.LoadServices(context.Background()); err != nil {
		defaultLogger.Fatal().Msgf(dictionary.GetInternalMessage(dictionary.ErrorLoadingService), config.ApplicationName, err.Error())
		return
	}

	// @mrz New Relic is ready at this point
	txn := services.NewRelic.StartTransaction("load_server")

	// (debugging: show services that are enabled or not)
	if appConfig.Debug {
		services.Logger.Debug().Msgf(
			"datastore: %s | cachestore: %s | taskmanager: %s | new_relic: %t | paymail: %t | graphql: %t",
			appConfig.Db.Datastore.Engine.String(),
			appConfig.Cache.Engine.String(),
			appConfig.TaskManager.Factory.String(),
			appConfig.NewRelic.Enabled,
			appConfig.Paymail.Enabled,
			appConfig.GraphQL.Enabled,
		)
	}

	// Create a new app server
	appServer := server.NewServer(appConfig, services)

	idleConnectionsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// We received an interrupt signal, shut down.
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err = appServer.Shutdown(ctx); err != nil {
			services.Logger.Fatal().Msgf("error shutting down: %s", err.Error())
		}

		close(idleConnectionsClosed)
	}()

	// End new relic txn
	txn.End()

	// Listen and serve
	services.Logger.Debug().Msgf("starting %s server version %s at port %d...", config.ApplicationName, config.Version, appConfig.Server.Port)
	appServer.Serve()

	<-idleConnectionsClosed
}
