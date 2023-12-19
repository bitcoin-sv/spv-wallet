package actions

import (
	"net/http"
	"time"

	"github.com/BuxOrg/bux"
	"github.com/BuxOrg/bux-server/config"
	"github.com/BuxOrg/bux-server/dictionary"
	"github.com/gofrs/uuid"
	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
	"github.com/mrz1836/go-parameters"
)

// Action is the configuration for the actions and related services
type Action struct {
	AppConfig *config.AppConfig
	Services  *config.AppServices
}

// NewStack is used for registering routes
func NewStack(appConfig *config.AppConfig,
	services *config.AppServices) (Action, *apirouter.InternalStack) {
	return Action{AppConfig: appConfig, Services: services}, apirouter.NewStack()
}

// RequireAuthentication checks and requires authentication for the related method
func (a *Action) RequireAuthentication(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {

		// Check the authentication
		var knownErr dictionary.ErrorMessage
		if req, knownErr = CheckAuthentication(a.AppConfig, a.Services.Bux, req, false, true); knownErr.Code > 0 {
			ReturnErrorResponse(w, req, knownErr, "")
			return
		}

		// Continue to next method
		fn(w, req, p)
	}
}

// RequireBasicAuthentication checks and requires authentication for the related method, but does not require signing
func (a *Action) RequireBasicAuthentication(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {

		// Check the authentication
		var knownErr dictionary.ErrorMessage
		if req, knownErr = CheckAuthentication(a.AppConfig, a.Services.Bux, req, false, false); knownErr.Code > 0 {
			ReturnErrorResponse(w, req, knownErr, "")
			return
		}

		// Continue to next method
		fn(w, req, p)
	}
}

// RequireAdminAuthentication checks and requires ADMIN authentication for the related method
func (a *Action) RequireAdminAuthentication(fn httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {

		// Check the authentication
		var knownErr dictionary.ErrorMessage
		if req, knownErr = CheckAuthentication(a.AppConfig, a.Services.Bux, req, true, true); knownErr.Code > 0 {
			ReturnErrorResponse(w, req, knownErr, "")
			return
		}

		// Continue to next method
		fn(w, req, p)
	}
}

// Request will process the request in the router
func (a *Action) Request(_ *apirouter.Router, h httprouter.Handle) httprouter.Handle {
	return Request(h, a)
}

// CheckAuthentication will check the authentication
func CheckAuthentication(appConfig *config.AppConfig, bux bux.ClientInterface, req *http.Request,
	adminRequired bool, requireSigning bool) (*http.Request, dictionary.ErrorMessage) {

	// Bad/Unknown scheme
	if appConfig.Authentication.Scheme != config.AuthenticationSchemeXpub {
		return req, dictionary.GetError(dictionary.ErrorAuthenticationScheme, appConfig.Authentication.Scheme)
	}

	// AuthenticateFromRequest using the xPub scheme
	var err error
	if req, err = bux.AuthenticateRequest(
		req.Context(),
		req, []string{appConfig.Authentication.AdminKey}, adminRequired,
		requireSigning && appConfig.Authentication.RequireSigning,
		appConfig.Authentication.SigningDisabled,
	); err != nil {
		return req, dictionary.GetError(dictionary.ErrorAuthenticationError, err.Error())
	}

	// Return an empty error message
	return req, dictionary.ErrorMessage{}
}

// Request will write the request to the logs before and after calling the handler
func Request(h httprouter.Handle, a *Action) httprouter.Handle {
	return parameters.MakeHTTPRouterParsedReq(func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		// Start the custom response writer
		guid, _ := uuid.NewV4()
		writer := &apirouter.APIResponseWriter{
			IPAddress:      apirouter.GetClientIPAddress(req),
			Method:         req.Method,
			RequestID:      guid.String(),
			ResponseWriter: w,
			Status:         0, // future use with E-tags
			URL:            req.URL.String(),
			UserAgent:      req.UserAgent(),
		}

		// Start the log (timer)
		start := time.Now()

		// Fire the request
		h(writer, req, ps)

		// Complete the timer and final log
		elapsed := time.Since(start)

		if a.AppConfig.RequestLogging {
			a.Services.Logger.Debug().Msgf("%d | %s | %s | %s | %s ", writer.Status, elapsed, req.RemoteAddr, req.Method, req.URL)
		}
	})
}
