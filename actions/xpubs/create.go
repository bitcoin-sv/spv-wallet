package xpubs

import (
	"net/http"

	"github.com/BuxOrg/bux"
	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// create will make a new model using the services defined in the action object
func (a *Action) create(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	// Parse the params
	params := apirouter.GetParams(req)

	// params
	key := params.GetString("key")
	metadata := params.GetJSON("metadata")

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
