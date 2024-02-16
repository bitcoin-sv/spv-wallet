package accesskeys

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/actions"
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// count will fetch a count of access keys filtered by metadata
// Count of access keys godoc
// @Summary		Count of access keys
// @Description	Count of access keys
// @Tags		Access-key
// @Produce		json
// @Param		metadata query string false "metadata"
// @Param		conditions query string false "conditions"
// @Success		200
// @Router		/v1/access-key/count [post]
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

	// Record a new transaction (get the hex from parameters)a
	var count int64
	if count, err = a.Services.SPV.GetAccessKeysByXPubIDCount(
		req.Context(),
		reqXPubID,
		metadata,
		conditions,
	); err != nil {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, err.Error())
		return
	}

	// Return response
	apirouter.ReturnResponse(w, req, http.StatusOK, count)
}
