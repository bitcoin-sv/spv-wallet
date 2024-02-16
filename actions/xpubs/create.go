package xpubs

import (
	"net/http"

	"github.com/bitcoin-sv/bux"
	"github.com/bitcoin-sv/spv-wallet/actions"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// create will make a new model using the services defined in the action object
// Create xPub godoc
// @Summary		Create xPub
// @Description	Create xPub
// @Tags		xPub
// @Produce		json
// @Param		key query string true "key"
// @Param		metadata query string false "metadata"
// @Success		201
// @Router		/v1/xpub [post]
// @Security	spv-wallet-auth-xpub
func (a *Action) create(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	// Parse the params
	params := apirouter.GetParams(req)

	// params
	key := params.GetString("key")
	metadata := params.GetJSON(actions.MetadataField)

	// Create a new xPub
	xPub, err := a.Services.SPV.NewXpub(
		req.Context(), key,
		bux.WithMetadatas(metadata),
	)
	if err != nil {
		apirouter.ReturnResponse(w, req, http.StatusUnprocessableEntity, err.Error())
		return
	}

	contract := mappings.MapToXpubContract(xPub)

	// Return response
	apirouter.ReturnResponse(w, req, http.StatusCreated, contract)
}
