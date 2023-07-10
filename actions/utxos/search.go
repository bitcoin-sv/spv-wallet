package utxos

import (
	"net/http"

	"github.com/BuxOrg/bux"
	buxmodels "github.com/BuxOrg/bux-models"
	"github.com/BuxOrg/bux-server/actions"
	"github.com/BuxOrg/bux-server/mappings"
	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// search will fetch a list of utxos filtered on conditions and metadata
// Search UTXO godoc
// @Summary		Search UTXO
// @Description	Search UTXO
// @Tags		UTXO
// @Produce		json
// @Param		page query int false "page"
// @Param		page_size query int false "page_size"
// @Param		order_by_field query string false "order_by_field"
// @Param		sort_direction query string false "sort_direction"
// @Param		metadata query string false "metadata"
// @Param		conditions query string false "conditions"
// @Success		200
// @Router		/v1/utxo/search [post]
// @Security	bux-auth-xpub
func (a *Action) search(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	reqXPubID, _ := bux.GetXpubIDFromRequest(req)

	// Parse the params
	params := apirouter.GetParams(req)
	queryParams, modelMetadata, conditions, err := actions.GetQueryParameters(params)
	metadata := mappings.MapToBuxMetadata(modelMetadata)
	if err != nil {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, err.Error())
		return
	}

	// Record a new transaction (get the hex from parameters)a
	var utxos []*bux.Utxo
	if utxos, err = a.Services.Bux.GetUtxosByXpubID(
		req.Context(),
		reqXPubID,
		metadata,
		conditions,
		queryParams,
	); err != nil {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, err.Error())
		return
	}

	contracts := make([]*buxmodels.Utxo, 0)
	for _, utxo := range utxos {
		contracts = append(contracts, mappings.MapToUtxoContract(utxo))
	}

	// Return response
	apirouter.ReturnResponse(w, req, http.StatusOK, bux.DisplayModels(contracts))
}
