/*
Package main is the core service layer for the xAPI service
*/
package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	pmail "github.com/BuxOrg/bux-server/actions/paymail"
	"github.com/BuxOrg/bux-server/config"
	"github.com/BuxOrg/bux-server/dictionary"
	"github.com/BuxOrg/bux-server/server"
	"github.com/mrz1836/go-logger"
)

// main method starts everything for the xAPI service
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

	// Set the Paymail Server Configuration (if enabled)
	if err = services.SetPaymailServer(
		appConfig,
		pmail.NewServiceProvider(services.Bux, appConfig),
	); err != nil {
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
			fmt.Sprintf("new_relic service: %t",
				appConfig.NewRelic.Enabled,
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
