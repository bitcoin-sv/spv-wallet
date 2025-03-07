// Package server is for all the SPV Wallet settings and HTTP server
package server

import (
	"context"
	"crypto/tls"
	"net/http"
	"strconv"

	"github.com/bitcoin-sv/spv-wallet/actions"
	"github.com/bitcoin-sv/spv-wallet/actions/paymailserver"
	v2 "github.com/bitcoin-sv/spv-wallet/actions/v2"
	"github.com/bitcoin-sv/spv-wallet/api"
	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/logging"
	"github.com/bitcoin-sv/spv-wallet/metrics"
	"github.com/bitcoin-sv/spv-wallet/server/handlers"
	"github.com/bitcoin-sv/spv-wallet/server/middleware"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
)

// Server is the configuration, services, and actual web server
type Server struct {
	AppConfig       *config.AppConfig
	Router          *gin.Engine
	SpvWalletEngine engine.ClientInterface
	WebServer       *http.Server
	Logger          zerolog.Logger
}

// NewServer will return a new server service
func NewServer(appConfig *config.AppConfig, spvWalletEngine engine.ClientInterface, logger zerolog.Logger) *Server {
	return &Server{
		AppConfig:       appConfig,
		SpvWalletEngine: spvWalletEngine,
		Logger:          logger,
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

	s.Logger.Debug().Msgf("starting %s server at port %d...", s.AppConfig.GetUserAgent(), s.AppConfig.Server.Port)
	// Listen and serve
	if err := s.WebServer.ListenAndServe(); err != nil {
		s.Logger.Info().Err(err).Msgf("shutting down %s server [%s] on port %d...", s.AppConfig.GetUserAgent(), err.Error(), s.AppConfig.Server.Port)
	}
}

// Shutdown will stop the web server
func (s *Server) Shutdown(ctx context.Context) error {
	err := s.WebServer.Shutdown(ctx)
	if err != nil {
		err = spverrors.Wrapf(err, "error shutting down server")
	}
	return err
}

// Handlers will return handlers
func (s *Server) Handlers() *gin.Engine {
	httpLogger := s.Logger.With().Str("service", "http-server").Logger()
	if httpLogger.GetLevel() > zerolog.DebugLevel {
		gin.SetMode(gin.ReleaseMode)
	}
	logging.SetGinWriters(&httpLogger)
	ginEngine := gin.New()
	ginEngine.Use(logging.GinMiddleware(httpLogger), gin.Recovery())
	ginEngine.Use(middleware.AppContextMiddleware(s.AppConfig, s.SpvWalletEngine, s.Logger))
	ginEngine.Use(middleware.CorsMiddleware())

	metrics.SetupGin(ginEngine)

	ginEngine.NoRoute(metrics.NoRoute, NotFound)
	ginEngine.NoMethod(MethodNotAllowed)

	s.Router = ginEngine

	setupServerRoutes(s.AppConfig, s.SpvWalletEngine, s.Router, &httpLogger)

	return s.Router
}

func setupServerRoutes(appConfig *config.AppConfig, spvWalletEngine engine.ClientInterface, ginEngine *gin.Engine, log *zerolog.Logger) {
	handlersManager := handlers.NewManager(ginEngine, appConfig)

	if !appConfig.ExperimentalFeatures.V2 {
		paymailserver.Register(spvWalletEngine.GetPaymailConfig().Configuration, ginEngine)
		actions.Register(handlersManager)
	} else {
		v2.RegisterNonOpenAPIRoutes(ginEngine, appConfig, spvWalletEngine)
		api.RegisterHandlersWithOptions(ginEngine, v2.NewV2API(appConfig, spvWalletEngine, log), api.GinServerOptions{
			BaseURL: "",
			Middlewares: []api.MiddlewareFunc{
				middleware.SignatureAuthWithScopes(log),
			},
			ErrorHandler: func(c *gin.Context, err error, statusCode int) {
				spverrors.ErrorResponse(c, err, log)
			},
		})
	}

	if appConfig.DebugProfiling {
		pprof.Register(ginEngine, "debug/pprof")
	}
}
