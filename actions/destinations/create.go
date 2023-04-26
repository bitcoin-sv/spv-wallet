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
// Count Destinations godoc
// @Summary      Create a new destination
// @Description  Create a new destination
// @Tags		 Destinations
// @Produce      json
// @Param 		 type query string false "type"
// @Param 		 reference_id query string false "reference_id"
// @Param 		 metadata query string false "metadata"
// @Success      201
// @Router       /v1/destination [post]
// @Security bux-auth-xpub
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

	opts := a.Services.Bux.DefaultModelOptions()

	if metadata != nil {
		opts = append(opts, bux.WithMetadatas(metadata))
	}

	// Get a new destination
	var destination *bux.Destination
	if destination, err = a.Services.Bux.NewDestination(
		req.Context(),
		xPub.RawXpub(),
		uint32(0), // todo: use a constant? protect this?
		scriptType,
		true, // monitor this address as it was created by request of a user to share
		opts...,
	); err != nil {
		apirouter.ReturnResponse(w, req, http.StatusUnprocessableEntity, err.Error())
		return
	}

	// Return response
	apirouter.ReturnResponse(w, req, http.StatusCreated, bux.DisplayModels(destination))
}
