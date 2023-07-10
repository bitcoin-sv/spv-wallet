package pmail

import (
	"net/http"

	"github.com/BuxOrg/bux"
	"github.com/BuxOrg/bux-server/mappings"
	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// create will make a new model using the services defined in the action object
// Create Paymail godoc
// @Summary		Create paymail
// @Description	Create paymail
// @Tags		Paymails
// @Param		key query string true "key"
// @Param		address query string true "address"
// @Param		public_name query string false "public_name"
// @Param		avatar query string false "avatar"
// @Param		metadata query string false "metadata"
// @Produce		json
// @Success		201
// @Router		/v1/paymail [post]
// @Security	bux-auth-xpub
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

	contract := mappings.MapToPaymailContract(paymailAddress)

	// Return response
	apirouter.ReturnResponse(w, req, http.StatusCreated, bux.DisplayModels(contract))
}
