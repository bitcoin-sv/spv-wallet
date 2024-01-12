package base

import (
	"net/http"
	"net/http/pprof"

	"github.com/BuxOrg/bux-server/actions"
	"github.com/BuxOrg/bux-server/config"
	apirouter "github.com/mrz1836/go-api-router"
)

// Action is an extension of actions.Action for this package
type Action struct {
	actions.Action
}

// RegisterRoutes register all the package specific routes
func RegisterRoutes(router *apirouter.Router, appConfig *config.AppConfig, services *config.AppServices) {
	// Load the actions and set the services
	action := &Action{actions.Action{AppConfig: appConfig, Services: services}}

	// Set the main index page (navigating to slash)
	router.HTTPRouter.GET("/", action.Request(router, router.Request(index)))
	router.HTTPRouter.OPTIONS("/", router.SetCrossOriginHeaders)
	router.HTTPRouter.HEAD("/", actions.Head)

	// Set the health request (used for load balancers)
	router.HTTPRouter.GET("/"+config.HealthRequestPath, router.RequestNoLogging(actions.Health))
	router.HTTPRouter.OPTIONS("/"+config.HealthRequestPath, router.SetCrossOriginHeaders)
	router.HTTPRouter.HEAD("/"+config.HealthRequestPath, router.SetCrossOriginHeaders)

	// Debugging (shows all the Go profiler information)
	if action.AppConfig.DebugProfiling {
		router.HTTPRouter.HandlerFunc(http.MethodGet, "/debug/pprof/", pprof.Index)
		router.HTTPRouter.HandlerFunc(http.MethodGet, "/debug/pprof/cmdline", pprof.Cmdline)
		router.HTTPRouter.HandlerFunc(http.MethodGet, "/debug/pprof/profile", pprof.Profile)
		router.HTTPRouter.HandlerFunc(http.MethodGet, "/debug/pprof/symbol", pprof.Symbol)
		router.HTTPRouter.HandlerFunc(http.MethodGet, "/debug/pprof/trace", pprof.Trace)
		router.HTTPRouter.Handler(http.MethodGet, "/debug/pprof/goroutine", pprof.Handler("goroutine"))
		router.HTTPRouter.Handler(http.MethodGet, "/debug/pprof/heap", pprof.Handler("heap"))
		router.HTTPRouter.Handler(http.MethodGet, "/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
		router.HTTPRouter.Handler(http.MethodGet, "/debug/pprof/block", pprof.Handler("block"))
	}

	// Set the 404 handler (any request not detected)
	router.HTTPRouter.NotFound = http.HandlerFunc(actions.NotFound)

	// Set the method not allowed
	router.HTTPRouter.MethodNotAllowed = http.HandlerFunc(actions.MethodNotAllowed)
}
