package pmail

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
	key := params.GetString("key")                // the rawXPubKey
	address := params.GetString("address")        // the full paymail address
	publicName := params.GetString("public_name") // the public name
	avatar := params.GetString("avatar")          // the avatar
	metadata := params.GetJSON("metadata")        // optional metadata

	opts := a.Services.Bux.DefaultModelOptions()

	if metadata != nil {
		opts = append(opts, bux.WithMetadatas(metadata))
	}

	paymailAddress, err := a.Services.Bux.NewPaymailAddress(req.Context(), key, address, publicName, avatar, opts...)
	if err != nil {
		apirouter.ReturnResponse(w, req, http.StatusUnprocessableEntity, err.Error())
		return
	}

	// Return response
	apirouter.ReturnResponse(w, req, http.StatusCreated, bux.DisplayModels(paymailAddress))
}
