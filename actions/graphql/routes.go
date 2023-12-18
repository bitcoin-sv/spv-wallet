package graphql

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"regexp"

	"github.com/99designs/gqlgen/graphql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/BuxOrg/bux"
	"github.com/BuxOrg/bux-server/actions"
	"github.com/BuxOrg/bux-server/config"
	"github.com/BuxOrg/bux-server/dictionary"
	"github.com/BuxOrg/bux-server/graph"
	"github.com/BuxOrg/bux-server/graph/generated"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

const (
	allowCredentialsHeader string = "Access-Control-Allow-Credentials"
	allowHeadersHeader     string = "Access-Control-Allow-Headers"
	allowMethodsHeader     string = "Access-Control-Allow-Methods"
	allowOriginHeader      string = "Access-Control-Allow-Origin"
)

type requestInfo struct {
	id        uuid.UUID
	method    string
	path      string
	ip        string
	userAgent string
}

// RegisterRoutes register all the package specific routes
func RegisterRoutes(router *apirouter.Router, appConfig *config.AppConfig, services *config.AppServices) {

	// Use the authentication middleware wrapper
	a, require := actions.NewStack(appConfig, services)
	// require.Use(a.RequireAuthentication)

	srv := handler.NewDefaultServer(generated.NewExecutableSchema(generated.Config{Resolvers: &graph.Resolver{}}))
	if appConfig.RequestLogging {
		re := regexp.MustCompile(`[\r?\n|\s+]`)
		srv.AroundOperations(func(ctx context.Context, next graphql.OperationHandler) graphql.ResponseHandler {
			oc := graphql.GetOperationContext(ctx)
			reqInfo := ctx.Value(config.GraphRequestInfo).(requestInfo)
			params := map[string]interface{}{
				"query":     re.ReplaceAllString(oc.RawQuery, " "),
				"variables": oc.Variables,
			}
			// LogParamsFormat "request_id=\"%s\" method=\"%s\" path=\"%s\" ip_address=\"%s\" user_agent=\"%s\" params=\"%v\"\n"
			services.Logger.Info().Msgf(apirouter.LogParamsFormat, reqInfo.id, reqInfo.method, reqInfo.path, reqInfo.ip, reqInfo.userAgent, params)
			return next(ctx)
		})
		srv.SetErrorPresenter(func(ctx context.Context, err error) *gqlerror.Error {
			// LogErrorFormat "request_id=\"%s\" ip_address=\"%s\" type=\"%s\" internal_message=\"%s\" code=%d\n"
			reqInfo := ctx.Value(config.GraphRequestInfo).(requestInfo)
			services.Logger.Info().Msgf(apirouter.LogErrorFormat, reqInfo.id, reqInfo.ip, "GraphQL", err.Error(), 500)
			return &gqlerror.Error{
				Message: "presented: " + err.Error(),
				Path:    graphql.GetPath(ctx),
			}
		})
	}

	// Set the handle
	h := require.Wrap(wrapHandler(
		router,
		a.AppConfig,
		a.Services,
		srv,
		true,
	))

	// Set the GET routes
	router.HTTPRouter.GET(serverPath, h)

	// Set the POST routes
	router.HTTPRouter.POST(serverPath, h)

	// only show in debug mode
	if appConfig.Debug {
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
				services.Logger.Debug().Msgf("started graphql playground server on %s", playgroundPath)
			}
		} else {
			services.Logger.Error().Msgf("Failed starting graphql playground server directory equals playground directory %s = %s", serverPath, playgroundPath)
		}
	}

	// Success on the routes
	if appConfig.Debug {
		services.Logger.Debug().Msg("registered graphql routes on " + serverPath)
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

		var err error
		if withAuth {
			var knownErr dictionary.ErrorMessage
			// the graphql cannot check the signature using the bux check, will add extra checks in the endpoints
			if req, knownErr = actions.CheckAuthentication(appConfig, services.Bux, req, false, false); knownErr.Code > 0 {
				err = errors.New(knownErr.PublicMessage)
			}
		}

		signed := req.Context().Value(bux.ParamAuthSigned)
		if signed == nil {
			signed = false
		}

		xPubID := req.Context().Value(bux.ParamXPubHashKey)
		if xPubID == nil {
			xPubID = ""
		}

		// Create the context
		ctx := context.WithValue(req.Context(), config.GraphConfigKey, &graph.GQLConfig{
			AppConfig: appConfig,
			Services:  services,
			Signed:    signed.(bool),
			XPub:      req.Header.Get(bux.AuthHeader),
			XPubID:    xPubID.(string),
			AuthError: err,
		})

		guid, _ := uuid.NewV4()
		ctx = context.WithValue(ctx, config.GraphRequestInfo, requestInfo{
			id:        guid,
			method:    req.Method,
			path:      req.RequestURI,
			ip:        req.RemoteAddr,
			userAgent: req.UserAgent(),
		})

		// Call your original http.Handler
		h.ServeHTTP(w, req.WithContext(ctx))
	}
}
