package actions

import (
	"net/http"

	"github.com/BuxOrg/bux-server/dictionary"
	apirouter "github.com/mrz1836/go-api-router"
)

// ReturnErrorResponse will return a response using a dictionary.
// Error (using standard error responses)
//
// This wrapper can be removed, but was added in case we wanted to change the behavior
// of the error response on the action level
func ReturnErrorResponse(w http.ResponseWriter, req *http.Request,
	error dictionary.ErrorMessage, data interface{}) {
	apirouter.ReturnResponse(
		w,
		req,
		error.StatusCode,
		apirouter.ErrorFromRequest(
			req, error.InternalMessage, error.PublicMessage,
			int(error.Code), error.StatusCode, data,
		),
	)
}
