package utxos

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/actions"
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// count will count all the utxos fulfilling the given conditions
// Count of UTXOs godoc
// @Summary		Count of UTXOs
// @Description	Count of UTXOs
// @Tags		UTXO
// @Produce		json
// @Param		metadata query string false "metadata"
// @Param		conditions query string false "conditions"
// @Success		200
// @Router		/v1/utxo/count [post]
// @Security	x-auth-xpub
func (a *Action) count(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	reqXPubID, _ := engine.GetXpubIDFromRequest(req)

	// Parse the params
	params := apirouter.GetParams(req)
	_, metadataModel, conditions, err := actions.GetQueryParameters(params)
	metadata := mappings.MapToSPVMetadata(metadataModel)
	if err != nil {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, err.Error())
		return
	}

	dbConditions := map[string]interface{}{}
	if conditions != nil {
		dbConditions = *conditions
	}
	// force the xpub_id of the current user on query
	dbConditions["xpub_id"] = reqXPubID

	// Get a utxo using a xPub
	var count int64
	if count, err = a.Services.SPV.GetUtxosCount(
		req.Context(),
		metadata,
		&dbConditions,
	); err != nil {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, err.Error())
		return
	}

	// Return response
	apirouter.ReturnResponse(w, req, http.StatusOK, count)
}
