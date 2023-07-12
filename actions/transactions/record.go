package transactions

import (
	"net/http"

	"github.com/BuxOrg/bux"
	"github.com/BuxOrg/bux-server/actions"
	"github.com/BuxOrg/bux-server/mappings"
	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// record will save and complete a transaction
// Record transaction godoc
// @Summary		Record transaction
// @Description	Record transaction
// @Tags		Transactions
// @Produce		json
// @Param		hex query string true "hex"
// @Param		reference_id query string true "reference_id"
// @Param		metadata query string false "metadata"
// @Success		200
// @Router		/v1/transaction/record [post]
// @Security	bux-auth-xpub
func (a *Action) record(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	// Parse the params
	params := apirouter.GetParams(req)

	// Get the xPub from the request (via authentication)
	reqXPub, _ := bux.GetXpubFromRequest(req)
	xPub, err := a.Services.Bux.GetXpub(req.Context(), reqXPub)
	if err != nil {
		apirouter.ReturnResponse(w, req, http.StatusUnprocessableEntity, err.Error())
		return
	} else if xPub == nil {
		apirouter.ReturnResponse(w, req, http.StatusForbidden, actions.ErrXpubNotFound)
		return
	}

	// Set the metadata
	metadata, ok := params.GetJSONOk(bux.ModelMetadata.String())
	opts := make([]bux.ModelOps, 0)
	if ok {
		opts = append(opts, bux.WithMetadatas(metadata))
	}

	// Record a new transaction (get the hex from parameters)
	var transaction *bux.Transaction
	if transaction, err = a.Services.Bux.RecordTransaction(
		req.Context(),
		reqXPub,
		params.GetString("hex"),
		params.GetString(bux.ReferenceIDField),
		opts...,
	); err != nil {
		apirouter.ReturnResponse(w, req, http.StatusUnprocessableEntity, err.Error())
		return
	}

	contract := mappings.MapToTransactionContract(transaction)

	// Return response
	apirouter.ReturnResponse(w, req, http.StatusCreated, contract)
}
