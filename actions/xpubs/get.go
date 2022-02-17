package xpubs

import (
	"net/http"

	"github.com/BuxOrg/bux"
	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// get will get an existing model
func (a *Action) get(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	reqXPub, _ := bux.GetXpubFromRequest(req)

	// Parse the params
	params := apirouter.GetParams(req)
	key := params.GetString("key")
	if key != "" {
		if isAdmin, ok := bux.IsAdminRequest(req); !isAdmin || !ok {
			apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, bux.ErrNotAdminKey)
			return
		}
	} else {
		key = reqXPub
	}

	// Get an xPub
	xPub, err := a.Services.Bux.GetXpub(
		req.Context(), key,
	)
	if err != nil {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, err.Error())
		return
	}

	signed := req.Context().Value("auth_signed")
	if signed == nil || !signed.(bool) {
		// remove private data from the returned xPub
		xPub.NextExternalNum = 0
		xPub.NextInternalNum = 0
		xPub.Metadata = nil
	}

	// Return response
	apirouter.ReturnResponse(w, req, http.StatusOK, bux.DisplayModels(xPub))
}
