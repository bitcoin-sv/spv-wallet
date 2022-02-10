package actions

import (
	"net/http"

	"github.com/BuxOrg/bux"
	"github.com/BuxOrg/bux-server/config"
	"github.com/BuxOrg/bux-server/dictionary"
	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
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
		if req, knownErr = CheckAuthentication(a.AppConfig, a.Services.Bux, req, false); knownErr.Code > 0 {
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
		if req, knownErr = CheckAuthentication(a.AppConfig, a.Services.Bux, req, true); knownErr.Code > 0 {
			ReturnErrorResponse(w, req, knownErr, "")
			return
		}

		// Continue to next method
		fn(w, req, p)
	}
}

// CheckAuthentication will check the authentication
func CheckAuthentication(appConfig *config.AppConfig, bux bux.ClientInterface, req *http.Request, adminRequired bool) (*http.Request, dictionary.ErrorMessage) {

	// Bad/Unknown scheme
	if appConfig.Authentication.Scheme != config.AuthenticationSchemeXpub {
		return req, dictionary.GetError(dictionary.ErrorAuthenticationScheme, appConfig.Authentication.Scheme)
	}

	// AuthenticateFromRequest using the xPub scheme
	var err error
	if req, err = bux.AuthenticateRequest(
		req.Context(),
		req, []string{appConfig.Authentication.AdminKey}, adminRequired,
		appConfig.Authentication.RequireSigning,
		appConfig.Authentication.SigningDisabled,
	); err != nil {
		return req, dictionary.GetError(dictionary.ErrorAuthenticationError, err.Error())
	}

	// Return an empty error message
	return req, dictionary.ErrorMessage{}
}
