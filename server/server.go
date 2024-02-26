// Package server is for all the SPV Wallet settings and HTTP server
package server

import (
	"context"
	"crypto/tls"
	"github.com/bitcoin-sv/spv-wallet/actions"
	"github.com/bitcoin-sv/spv-wallet/actions/base"
	"github.com/bitcoin-sv/spv-wallet/logging"
	"github.com/bitcoin-sv/spv-wallet/metrics"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"net/http"
	"strconv"

	"errors"
	accesskeys "github.com/bitcoin-sv/spv-wallet/actions/access_keys"
	"github.com/bitcoin-sv/spv-wallet/actions/admin"
	"github.com/bitcoin-sv/spv-wallet/actions/destinations"
	"github.com/bitcoin-sv/spv-wallet/actions/transactions"
	"github.com/bitcoin-sv/spv-wallet/actions/utxos"
	"github.com/bitcoin-sv/spv-wallet/actions/xpubs"
	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
	router "github.com/bitcoin-sv/spv-wallet/server/routes"
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
	return s.WebServer.Shutdown(ctx)
}

// Handlers will return handlers
func (s *Server) Handlers() *gin.Engine {
	// Start a transaction for loading handlers
	txn := s.Services.NewRelic.StartTransaction("load_handlers")
	defer txn.End()

	segment := txn.StartSegment("create_router")

	httpLogger := s.Services.Logger.With().Str("service", "http-server").Logger()
	if httpLogger.GetLevel() > zerolog.DebugLevel {
		gin.SetMode(gin.ReleaseMode)
	}
	logging.SetGinWriters(&httpLogger)
	engine := gin.New()
	engine.Use(logging.GinMiddleware(&httpLogger), gin.Recovery())
	engine.Use(auth.CorsMiddleware())

	s.Router = engine

	segment.End()

	// Start the segment
	defer txn.StartSegment("register_handlers").End()

	SetupServerRoutes(s.AppConfig, s.Services, s.Router)

	return s.Router
}

func SetupServerRoutes(appConfig *config.AppConfig, services *config.AppServices, engine *gin.Engine) {
	adminRoutes := admin.NewHandler(appConfig, services)
	baseRoutes := base.NewHandler(appConfig, engine)

	accessKeyApiRoutes := accesskeys.NewHandler(appConfig, services)
	destinationBasicRoutes, destinationApiRoutes := destinations.NewHandler(appConfig, services)
	transactionBasicRoutes, transactionApiRoutes, transactionCallbackRoutes := transactions.NewHandler(appConfig, services)
	utxoApiRoutes := utxos.NewHandler(appConfig, services)
	xPubApiRoutes := xpubs.NewHandler(appConfig, services)

	routes := []interface{}{
		// Admin routes
		adminRoutes,
		// Base routes
		baseRoutes,
		// Access key routes
		accessKeyApiRoutes,
		// Destination routes
		destinationBasicRoutes,
		destinationApiRoutes,
		// Transaction routes
		transactionBasicRoutes,
		transactionApiRoutes,
		transactionCallbackRoutes,
		// Utxo routes
		utxoApiRoutes,
		// xPub routes
		xPubApiRoutes,
	}

	prefix := "/" + config.APIVersion
	baseRouter := engine.Group("")
	authRouter := engine.Group("", auth.BasicMiddleware(services.SpvWalletEngine, appConfig))
	basicAuthRouter := authRouter.Group(prefix, auth.SignatureMiddleware(appConfig, false, false))
	apiAuthRouter := authRouter.Group(prefix, auth.SignatureMiddleware(appConfig, true, false))
	adminAuthRouter := authRouter.Group(prefix, auth.SignatureMiddleware(appConfig, true, true), auth.AdminMiddleware())
	callbackAuthRouter := baseRouter.Group("", auth.CallbackTokenMiddleware(appConfig))

	for _, r := range routes {
		switch r := r.(type) {
		case router.AdminEndpoints:
			r.RegisterAdminEndpoints(adminAuthRouter)
		case router.APIEndpoints:
			r.RegisterAPIEndpoints(apiAuthRouter)
		case router.BasicEndpoints:
			r.RegisterBasicEndpoints(basicAuthRouter)
		case router.BaseEndpoints:
			r.RegisterBaseEndpoints(baseRouter)
		case router.CallbackEndpoints:
			r.RegisterCallbackEndpoints(callbackAuthRouter)
		default:
			panic(errors.New("unexpected router endpoints registrar"))
		}
	}

	// Register paymail routes
	services.SpvWalletEngine.GetPaymailConfig().RegisterRoutes(engine)

	// Set the 404 handler (any request not detected)
	engine.NoRoute(actions.NotFound)

	// Set the method not allowed
	engine.NoMethod(actions.MethodNotAllowed)

	registerSwaggerEndpoints(engine)

	if metrics, enabled := metrics.Get(); enabled {
		engine.GET("/metrics", gin.WrapH(metrics.HTTPHandler()))
	}
}
