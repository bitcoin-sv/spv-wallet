package xpubs

import (
	"net/http"

	"github.com/BuxOrg/bux"
	"github.com/BuxOrg/bux-server/actions"
	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// create will make a new model using the services defined in the action object
func (a *Action) createWithToken(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	// Parse the params
	params := apirouter.GetParams(req)

	// params
	key := params.GetString("key")
	token := params.GetString("token")
	metadata := params.GetJSON(actions.MetadataField)

	// TODO: this is temporary during testing, tokens still need to be fleshed out
	if token != "61zOLgJNf7q4apobdwRPfFAbGTD1zferPt-kpAexlGq" { //nolint:gosec // only temp
		apirouter.ReturnResponse(w, req, http.StatusUnprocessableEntity, "invalid token given")
		return
	}

	// Create a new xPub
	xPub, err := a.Services.Bux.NewXpub(
		req.Context(), key,
		bux.WithMetadatas(metadata),
	)
	if err != nil {
		apirouter.ReturnResponse(w, req, http.StatusUnprocessableEntity, err.Error())
		return
	}

	// Return response
	apirouter.ReturnResponse(w, req, http.StatusCreated, bux.DisplayModels(xPub))
}
