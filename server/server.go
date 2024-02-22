// Package server is for all the SPV Wallet settings and HTTP server
package server

import (
	"context"
	"crypto/tls"
	"net/http"
	"strconv"

	accesskeys "github.com/bitcoin-sv/spv-wallet/actions/access_keys"
	"github.com/bitcoin-sv/spv-wallet/actions/admin"
	"github.com/bitcoin-sv/spv-wallet/actions/base"
	"github.com/bitcoin-sv/spv-wallet/actions/contacts"
	"github.com/bitcoin-sv/spv-wallet/actions/destinations"
	pmail "github.com/bitcoin-sv/spv-wallet/actions/paymail"
	"github.com/bitcoin-sv/spv-wallet/actions/transactions"
	"github.com/bitcoin-sv/spv-wallet/actions/utxos"
	"github.com/bitcoin-sv/spv-wallet/actions/xpubs"
	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/bitcoin-sv/spv-wallet/metrics"
	apirouter "github.com/mrz1836/go-api-router"
	"github.com/newrelic/go-agent/v3/integrations/nrhttprouter"
	httpSwagger "github.com/swaggo/http-swagger"
)

// Server is the configuration, services, and actual web server
type Server struct {
	AppConfig *config.AppConfig
	Router    *apirouter.Router
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
func (s *Server) Handlers() *nrhttprouter.Router {
	// Start a transaction for loading handlers
	txn := s.Services.NewRelic.StartTransaction("load_handlers")
	defer txn.End()

	// Create a new router
	segment := txn.StartSegment("create_router")
	s.Router = apirouter.NewWithNewRelic(s.Services.NewRelic)
	s.Router.HTTPRouter.Handler(http.MethodGet, "/swagger", http.RedirectHandler("/swagger/index.html", http.StatusMovedPermanently))
	s.Router.HTTPRouter.Handler(http.MethodGet, "/swagger/*any", httpSwagger.WrapHandler)
	if metrics, enabled := metrics.Get(); enabled {
		s.Router.HTTPRouter.Handler(http.MethodGet, "/metrics", metrics.HTTPHandler())
	}
	segment.End()

	// Turned on all CORs - should be able to access in a browser
	s.Router.CrossOriginEnabled = true
	s.Router.CrossOriginAllowCredentials = true
	s.Router.CrossOriginAllowOrigin = "*"
	s.Router.CrossOriginAllowMethods = "POST,GET,OPTIONS,DELETE"
	s.Router.CrossOriginAllowHeaders = "*"

	// Start the segment
	defer txn.StartSegment("register_handlers").End()

	// Register all handlers (actions / routes)
	base.RegisterRoutes(s.Router, s.AppConfig, s.Services)
	admin.RegisterRoutes(s.Router, s.AppConfig, s.Services)
	accesskeys.RegisterRoutes(s.Router, s.AppConfig, s.Services)
	destinations.RegisterRoutes(s.Router, s.AppConfig, s.Services)
	transactions.RegisterRoutes(s.Router, s.AppConfig, s.Services)
	utxos.RegisterRoutes(s.Router, s.AppConfig, s.Services)
	xpubs.RegisterRoutes(s.Router, s.AppConfig, s.Services)
	contacts.RegisterRoutes(s.Router, s.AppConfig, s.Services)

	// Load Paymail
	pmail.RegisterRoutes(s.Router, s.AppConfig, s.Services)

	// Return the router
	return s.Router.HTTPRouter
}
