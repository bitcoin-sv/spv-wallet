/*
Package main is the core service layer for the SPV Wallet
*/
package main

import (
	"context"
	"os"
	"os/signal"
	"time"

	"github.com/bitcoin-sv/spv-wallet/config"
	_ "github.com/bitcoin-sv/spv-wallet/docs"
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/initializer"
	"github.com/bitcoin-sv/spv-wallet/logging"
	"github.com/bitcoin-sv/spv-wallet/server"
)

// version of the application that can be overridden with ldflags during build
// (e.g. go build -ldflags "-X main.version=1.2.3").
var version = "development"

// main method starts everything for the SPV Wallet
// @title           SPV Wallet
// @version         v1.0.0-beta
// @securityDefinitions.apikey x-auth-xpub
// @in header
// @name x-auth-xpub
// @securityDefinitions.apikey callback-auth
// @in header
// @name authorization
func main() {
	defaultLogger := logging.GetDefaultLogger()

	// Load the Application Configuration
	appConfig, err := config.Load(version, defaultLogger)
	if err != nil {
		defaultLogger.Fatal().Err(err).Msg("Error while loading configuration")
		return
	}

	// Validate configuration (before services have been loaded)
	if err = appConfig.Validate(); err != nil {
		defaultLogger.Fatal().Err(err).Msg("Invalid configuration")
		return
	}

	logger, err := logging.CreateLoggerWithConfig(appConfig)
	if err != nil {
		defaultLogger.Fatal().Err(err).Msg("Error while creating logger")
		return
	}

	appCtx := context.Background()

	opts, err := initializer.ToEngineOptions(appConfig, logger)
	if err != nil {
		defaultLogger.Fatal().Err(err).Msg("Error while creating engine options")
		return
	}

	spvWalletEngine, err := engine.NewClient(appCtx, opts...)
	if err != nil {
		defaultLogger.Fatal().Err(err).Msg("Error while creating SPV Wallet Engine")
		return
	}

	if appConfig.IsBeefEnabled() {
		spvWalletEngine.LogBHSReadiness(appCtx)
	}

	// Create a new app server
	appServer := server.NewServer(appConfig, spvWalletEngine, logger)

	idleConnectionsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// We received an interrupt signal, shut down.
		ctx, cancel := context.WithTimeout(appCtx, 5*time.Second)
		defer cancel()
		fatal := false
		if err = spvWalletEngine.Close(ctx); err != nil {
			logger.Error().Err(err).Msg("error when closing the engine")
			fatal = true
		}

		if err = appServer.Shutdown(ctx); err != nil {
			logger.Error().Err(err).Msg("error shutting down the server")
			fatal = true
		}

		close(idleConnectionsClosed)
		if fatal {
			os.Exit(1)
		}
	}()

	// Listen and serve
	appServer.Serve()

	<-idleConnectionsClosed
}
