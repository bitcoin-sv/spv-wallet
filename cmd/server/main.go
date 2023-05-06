/*
Package main is the core service layer for the BUX Server
*/
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/BuxOrg/bux-server/config"
	"github.com/BuxOrg/bux-server/dictionary"
	_ "github.com/BuxOrg/bux-server/docs"
	"github.com/BuxOrg/bux-server/server"
	"github.com/mrz1836/go-logger"
)

// main method starts everything for the BUX Server
// @title           BUX: Server
// @version         v0.5.0
// @securityDefinitions.apikey bux-auth-xpub
// @in header
// @name bux-auth-xpub
func main() {

	// Load the Application Configuration
	appConfig, err := config.Load("")
	if err != nil {
		logger.Fatalf(dictionary.GetInternalMessage(dictionary.ErrorLoadingConfig), err.Error())
		return
	}

	// Load the Application Services
	var services *config.AppServices
	if services, err = appConfig.LoadServices(context.Background()); err != nil {
		logger.Fatalf(dictionary.GetInternalMessage(dictionary.ErrorLoadingService), config.ApplicationName, err.Error())
		return
	}

	// @mrz New Relic is ready at this point
	txn := services.NewRelic.StartTransaction("load_server")

	// Validate configuration (after services have been loaded)
	if err = appConfig.Validate(txn); err != nil {
		logger.Fatalf(dictionary.GetInternalMessage(dictionary.ErrorLoadingConfig), err.Error())
		return
	}

	// (debugging: show services that are enabled or not)
	if appConfig.Debug {
		logger.Data(2, logger.DEBUG,
			fmt.Sprintf("datastore: %s | cachestore: %s | taskmanager: %s [%s] | new_relic: %t | paymail: %t | graphql: %t",
				appConfig.Datastore.Engine.String(),
				appConfig.Cachestore.Engine.String(),
				appConfig.TaskManager.Engine.String(),
				appConfig.TaskManager.Factory.String(),
				appConfig.NewRelic.Enabled,
				appConfig.Paymail.Enabled,
				appConfig.GraphQL.Enabled,
			),
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
			logger.Fatalf("error shutting down: %s", err.Error())
		}

		close(idleConnectionsClosed)
	}()

	// End new relic txn
	txn.End()

	// Listen and serve
	logger.Data(2, logger.DEBUG,
		"starting ["+appConfig.Environment+"] "+config.ApplicationName+" server...",
		logger.MakeParameter("port", appConfig.Server.Port),
	)
	appServer.Serve()

	<-idleConnectionsClosed
}
