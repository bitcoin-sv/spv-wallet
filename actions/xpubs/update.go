package xpubs

import (
	"net/http"

	"github.com/BuxOrg/bux"
	"github.com/BuxOrg/bux-server/actions"
	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// get will get an existing model
func (a *Action) update(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	reqXPub, _ := bux.GetXpubFromRequest(req)
	reqXPubID, _ := bux.GetXpubIDFromRequest(req)

	// Parse the params
	params := apirouter.GetParams(req)
	metadata := params.GetJSON(actions.MetadataField)

	// Get an xPub
	var xPub *bux.Xpub
	var err error
	xPub, err = a.Services.Bux.GetXpubByID(
		req.Context(), reqXPubID,
	)
	if err != nil {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, err.Error())
		return
	}

	if xPub.Metadata == nil {
		xPub.Metadata = make(bux.Metadata)
	}

	for key, value := range metadata {
		if value == nil {
			delete(xPub.Metadata, key)
		} else {
			xPub.Metadata[key] = value
		}
	}

	err = xPub.Save(req.Context())
	if err != nil {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, err.Error())
		return
	}

	signed := req.Context().Value("auth_signed")
	if signed == nil || !signed.(bool) || reqXPub == "" {
		xPub.RemovePrivateData()
	}

	// Return response
	apirouter.ReturnResponse(w, req, http.StatusOK, bux.DisplayModels(xPub))
}
