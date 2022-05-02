package admin

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// status will return the status of the admin login
func (a *Action) stats(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	stats, err := a.Services.Bux.GetStats(req.Context())
	if err != nil {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, err.Error())
		return
	}

	// Return response
	apirouter.ReturnResponse(w, req, http.StatusOK, stats)
}
