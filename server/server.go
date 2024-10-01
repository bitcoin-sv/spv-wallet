// Package server is for all the SPV Wallet settings and HTTP server
package server

import (
	"context"
	"crypto/tls"
	"net/http"
	"strconv"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"

	"github.com/bitcoin-sv/spv-wallet/actions"
	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/logging"
	"github.com/bitcoin-sv/spv-wallet/metrics"
	"github.com/bitcoin-sv/spv-wallet/server/handlers"
	"github.com/bitcoin-sv/spv-wallet/server/middleware"
)

// Server is the configuration, services, and actual web server
type Server struct {
	AppConfig *config.AppConfig
	Router    *gin.Engine
	Services  *config.AppServices
	WebServer *http.Server
}

// NewServer will return a new server service
func NewServer(appConfig *config.AppConfig, services *config.AppServices) *Server {
	return &Server{
		AppConfig: appConfig,
		Services:  services,
	}
}

// Serve will load a server and start serving
func (s *Server) Serve() {
	// Load the server defaults
	s.WebServer = &http.Server{
		Addr:              ":" + strconv.Itoa(s.AppConfig.Server.Port),
		Handler:           s.Handlers(),
		IdleTimeout:       s.AppConfig.Server.IdleTimeout,
		ReadTimeout:       s.AppConfig.Server.ReadTimeout,
		ReadHeaderTimeout: s.AppConfig.Server.ReadTimeout,
		WriteTimeout:      s.AppConfig.Server.WriteTimeout,
		TLSConfig: &tls.Config{
			NextProtos:       []string{"h2", "http/1.1"},
			MinVersion:       tls.VersionTLS12,
			CurvePreferences: []tls.CurveID{tls.CurveP256, tls.X25519},
			CipherSuites: []uint16{
				tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
				tls.TLS_ECDHE_ECDSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
				tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
				tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			},
			PreferServerCipherSuites: true,
		},
	}

	// Turn off keep alive
	// s.WebServer.SetKeepAlivesEnabled(false)

	// Listen and serve
	if err := s.WebServer.ListenAndServe(); err != nil {
		s.Services.Logger.Debug().Msgf("shutting down %s server [%s] on port %d...", config.ApplicationName, err.Error(), s.AppConfig.Server.Port)
	}
}

// Shutdown will stop the web server
func (s *Server) Shutdown(ctx context.Context) error {
	s.Services.CloseAll(ctx) // Should have been executed in main.go, but might panic and not run?
	err := s.WebServer.Shutdown(ctx)
	if err != nil {
		err = spverrors.Wrapf(err, "error shutting down server")
		return err
	}
	return nil
}

// Handlers will return handlers
func (s *Server) Handlers() *gin.Engine {

	httpLogger := s.Services.Logger.With().Str("service", "http-server").Logger()
	if httpLogger.GetLevel() > zerolog.DebugLevel {
		gin.SetMode(gin.ReleaseMode)
	}
	logging.SetGinWriters(&httpLogger)
	engine := gin.New()
	engine.Use(logging.GinMiddleware(&httpLogger), gin.Recovery())
	engine.Use(middleware.AppContextMiddleware(s.AppConfig, s.Services.SpvWalletEngine, s.Services.Logger))
	engine.Use(middleware.CorsMiddleware())

	metrics.SetupGin(engine)

	engine.NoRoute(metrics.NoRoute, NotFound)
	engine.NoMethod(MethodNotAllowed)

	s.Router = engine

	setupServerRoutes(s.AppConfig, s.Services, s.Router)

	return s.Router
}

func setupServerRoutes(appConfig *config.AppConfig, services *config.AppServices, ginEngine *gin.Engine) {
	handlersManager := handlers.NewManager(ginEngine, config.APIVersion)
	actions.Register(appConfig, handlersManager)

	services.SpvWalletEngine.GetPaymailConfig().RegisterRoutes(ginEngine)

	if appConfig.DebugProfiling {
		pprof.Register(ginEngine, "debug/pprof")
	}
}
