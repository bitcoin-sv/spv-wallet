package transactions

import (
	"net/http"

	"github.com/BuxOrg/bux"
	"github.com/BuxOrg/bux-server/mappings"
	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// get will fetch a transaction
// Get transaction by id godoc
// @Summary		Get transaction by id
// @Description	Get transaction by id
// @Tags		Transactions
// @Produce		json
// @Param		id query string true "id"
// @Success		200
// @Router		/v1/transaction [get]
// @Security	bux-auth-xpub
func (a *Action) get(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	// Parse the params
	params := apirouter.GetParams(req)

	// Get the xPub from the request (via authentication)
	reqXPubID, _ := bux.GetXpubIDFromRequest(req)

	// Get a transaction by ID
	transaction, err := a.Services.Bux.GetTransaction(
		req.Context(),
		reqXPubID,
		params.GetString("id"),
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
