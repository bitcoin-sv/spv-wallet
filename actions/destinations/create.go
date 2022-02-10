package destinations

import (
	"net/http"

	"github.com/BuxOrg/bux"
	"github.com/BuxOrg/bux-server/actions"
	"github.com/BuxOrg/bux/utils"
	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// create will make a new destination
func (a *Action) create(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	// Parse the params
	params := apirouter.GetParams(req)

	// Get the xPub from the request (via authentication)
	reqXPub, _ := bux.GetXpubFromRequest(req)
	xPub, err := a.Services.Bux.GetXpub(req.Context(), reqXPub)
	if err != nil {
		apirouter.ReturnResponse(w, req, http.StatusUnprocessableEntity, err.Error())
		return
	} else if xPub == nil {
		apirouter.ReturnResponse(w, req, http.StatusForbidden, actions.ErrXpubNotFound)
		return
	}

	// Get metadata (if any)
	metadata := params.GetJSON(bux.ModelMetadata.String())

	// Get the type
	scriptType := params.GetString("type")
	if scriptType == "" {
		scriptType = utils.ScriptTypePubKeyHash
	}

	// Set the reference ID
	referenceID := params.GetString(bux.ReferenceIDField)
	if len(referenceID) > 0 {
		metadata[bux.ReferenceIDField] = referenceID
	}

	// Get a new destination
	var destination *bux.Destination
	if destination, err = a.Services.Bux.NewDestination(
		req.Context(),
		xPub.RawXpub(),
		uint32(0), // todo: use a constant? protect this?
		scriptType,
		&metadata, // todo: check this exists before using a pointer?
	); err != nil {
		apirouter.ReturnResponse(w, req, http.StatusUnprocessableEntity, err.Error())
		return
	}

	// Return response
	apirouter.ReturnResponse(w, req, http.StatusCreated, bux.DisplayModels(destination))
}
