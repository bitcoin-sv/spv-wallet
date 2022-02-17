package graphql

import (
	"context"
	"fmt"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/BuxOrg/bux"
	"github.com/BuxOrg/bux-server/actions"
	"github.com/BuxOrg/bux-server/config"
	"github.com/BuxOrg/bux-server/dictionary"
	"github.com/BuxOrg/bux-server/graph"
	"github.com/BuxOrg/bux-server/graph/generated"
	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
	"github.com/mrz1836/go-logger"
)

const (
	allowCredentialsHeader string = "Access-Control-Allow-Credentials"
	allowHeadersHeader     string = "Access-Control-Allow-Headers"
	allowMethodsHeader     string = "Access-Control-Allow-Methods"
	allowOriginHeader      string = "Access-Control-Allow-Origin"
)

// RegisterRoutes register all the package specific routes
func RegisterRoutes(router *apirouter.Router, appConfig *config.AppConfig, services *config.AppServices) {

	// Use the authentication middleware wrapper
	a, require := actions.NewStack(appConfig, services)
	// require.Use(a.RequireAuthentication)

	// Set the path
	serverPath := appConfig.GraphQL.ServerPath
	if len(serverPath) == 0 {
		serverPath = defaultServerPath
	}

	// Set the handle
	h := require.Wrap(wrapHandler(
		router,
		a.AppConfig,
		a.Services,
		handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}})),
		true,
	))

	// Set the GET routes
	router.HTTPRouter.GET(serverPath, h)

	// Set the POST routes
	router.HTTPRouter.POST(serverPath, h)

	// only show in development mode
	if appConfig.Environment == config.EnvironmentDevelopment {
		playgroundPath := appConfig.GraphQL.PlaygroundPath
		if len(playgroundPath) == 0 {
			playgroundPath = defaultPlaygroundPath
		}
		if serverPath != playgroundPath {
			router.HTTPRouter.GET(
				playgroundPath,
				wrapHandler(
					router,
					appConfig,
					services,
					playground.Handler("GraphQL playground", serverPath),
					false,
				),
			)
			if appConfig.Debug {
				logger.Data(2, logger.DEBUG, "started graphql playground server on "+playgroundPath)
			}
		} else {
			logger.Data(2, logger.ERROR, "Failed starting graphql playground server directory equals playground directory "+serverPath+" = "+playgroundPath)
		}
	}

	// Success on the routes
	if appConfig.Debug {
		logger.Data(2, logger.DEBUG, "registered graphql routes on "+serverPath)
	}
}

// wrapHandler will wrap the "httprouter" with a generic handler
func wrapHandler(router *apirouter.Router, appConfig *config.AppConfig, services *config.AppServices,
	h http.Handler, withAuth bool) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {

		w.Header().Set(allowCredentialsHeader, fmt.Sprintf("%t", router.CrossOriginAllowCredentials))
		w.Header().Set(allowMethodsHeader, router.CrossOriginAllowMethods)
		w.Header().Set(allowHeadersHeader, router.CrossOriginAllowHeaders)
		w.Header().Set(allowOriginHeader, router.CrossOriginAllowOrigin)

		if withAuth {
			var knownErr dictionary.ErrorMessage
			// the graphql cannot check the signature using the bux check, will add extra checks in the endpoints
			if req, knownErr = actions.CheckAuthentication(appConfig, services.Bux, req, false, false); knownErr.Code > 0 {
				actions.ReturnErrorResponse(w, req, knownErr, "")
				return
			}
		}

		signed := req.Context().Value("auth_signed")
		if signed == nil {
			signed = false
		}

		// Create the context
		ctx := context.WithValue(req.Context(), config.GraphConfigKey, &graph.GQLConfig{
			AppConfig: appConfig,
			Services:  services,
			Signed:    signed.(bool),
			XPub:      req.Header.Get(bux.AuthHeader),
		})

		// Call your original http.Handler
		h.ServeHTTP(w, req.WithContext(ctx))
	}
}
