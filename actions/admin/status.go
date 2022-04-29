package admin

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// status will return the status of the admin login
func (a *Action) status(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	// Return response
	apirouter.ReturnResponse(w, req, http.StatusOK, true)
}
