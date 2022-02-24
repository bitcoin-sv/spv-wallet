package accesskeys

import (
	"net/http"

	"github.com/BuxOrg/bux/utils"

	"github.com/BuxOrg/bux"
	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// get will get an existing model
func (a *Action) get(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	reqXPub, _ := bux.GetXpubFromRequest(req)

	// Parse the params
	params := apirouter.GetParams(req)
	id := params.GetString("id")

	if id == "" {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, bux.ErrMissingFieldID)
		return
	}

	// Get an access key
	accessKey, err := a.Services.Bux.GetAccessKey(
		req.Context(), reqXPub, id,
	)
	if err != nil {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, err.Error())
		return
	}

	if accessKey.XpubID != utils.Hash(reqXPub) {
		apirouter.ReturnResponse(w, req, http.StatusForbidden, "unauthorized")
		return
	}

	// Return response
	apirouter.ReturnResponse(w, req, http.StatusOK, bux.DisplayModels(accessKey))
}
