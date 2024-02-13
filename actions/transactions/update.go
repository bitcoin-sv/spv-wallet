package transactions

import (
	"net/http"

	"github.com/BuxOrg/bux"
	"github.com/BuxOrg/spv-wallet/actions"
	"github.com/BuxOrg/spv-wallet/mappings"
	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// update will update a transaction
// Update transaction godoc
// @Summary		Update transaction
// @Description	Update transaction
// @Tags		Transactions
// @Produce		json
// @Param		id query string true "id"
// @Param		metadata query string true "metadata"
// @Success		200
// @Router		/v1/transaction [patch]
// @Security	spv-wallet-auth-xpub
func (a *Action) update(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	// Get the xPub from the request (via authentication)
	reqXPubID, _ := bux.GetXpubIDFromRequest(req)

	// Parse the params
	params := apirouter.GetParams(req)
	metadata := params.GetJSON(actions.MetadataField)

	// Get a transaction by ID
	transaction, err := a.Services.SPV.UpdateTransactionMetadata(
		req.Context(),
		reqXPubID,
		params.GetString("id"),
		metadata,
	)
	if err != nil {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, err.Error())
		return
	} else if transaction == nil {
		apirouter.ReturnResponse(w, req, http.StatusNotFound, "")
	} else if !transaction.IsXpubIDAssociated(reqXPubID) {
		apirouter.ReturnResponse(w, req, http.StatusForbidden, "unauthorized")
		return
	}

	contract := mappings.MapToTransactionContract(transaction)

	// Return response
	apirouter.ReturnResponse(w, req, http.StatusOK, contract)
}
