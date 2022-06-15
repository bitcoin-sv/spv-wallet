package admin

import (
	"net/http"

	"github.com/BuxOrg/bux"
	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// transactionRecord will save and complete a transaction directly, without any checks
func (a *Action) transactionRecord(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	// Parse the params
	params := apirouter.GetParams(req)

	// Set the metadata
	opts := make([]bux.ModelOps, 0)

	// Record a new transaction (get the hex from parameters)
	transaction, err := a.Services.Bux.RecordRawTransaction(
		req.Context(),
		params.GetString("hex"),
		opts...,
	)
	if err != nil {
		apirouter.ReturnResponse(w, req, http.StatusUnprocessableEntity, err.Error())
		return
	}

	// Return response
	apirouter.ReturnResponse(w, req, http.StatusCreated, bux.DisplayModels(transaction))
}
