package admin

import (
	"errors"
	"net/http"

	"github.com/BuxOrg/bux"
	"github.com/BuxOrg/bux-server/mappings"
	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
	"github.com/mrz1836/go-datastore"
)

// transactionRecord will save and complete a transaction directly, without any checks
// Record transactions godoc
// @Summary		Record transactions
// @Description	Record transactions
// @Tags		Admin
// @Produce		json
// @Param		hex query string true "Transaction hex"
// @Success		201
// @Router		/v1/admin/transactions/record [post]
// @Security	bux-auth-xpub
func (a *Action) transactionRecord(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	// Parse the params
	params := apirouter.GetParams(req)

	hex := params.GetString("hex")

	// Set the metadata
	opts := make([]bux.ModelOps, 0)

	// Record a new transaction (get the hex from parameters)
	transaction, err := a.Services.Bux.RecordRawTransaction(
		req.Context(),
		hex,
		opts...,
	)
	if err != nil {
		if errors.Is(err, datastore.ErrDuplicateKey) {
			// already registered, just return the registered transaction
			if transaction, err = a.Services.Bux.GetTransactionByHex(req.Context(), hex); err != nil {
				apirouter.ReturnResponse(w, req, http.StatusUnprocessableEntity, err.Error())
				return
			}
		} else {
			apirouter.ReturnResponse(w, req, http.StatusUnprocessableEntity, err.Error())
			return
		}
	}

	contract := mappings.MapToTransactionContract(transaction)

	// Return response
	apirouter.ReturnResponse(w, req, http.StatusCreated, bux.DisplayModels(contract))
}
