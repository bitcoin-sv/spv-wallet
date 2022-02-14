package transactions

import (
	"encoding/json"
	"net/http"

	"github.com/BuxOrg/bux"
	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// list will fetch a list of transactions
func (a *Action) list(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	reqXPub, _ := bux.GetXpubFromRequest(req)

	// Parse the params
	params := apirouter.GetParams(req)

	metadataReq := params.GetJSON("metadata")
	var metadata *bux.Metadata
	if len(metadataReq) > 0 {
		// marshal the metadata into the Metadata model
		metaJSON, _ := json.Marshal(metadataReq) // nolint: errchkjson // ignore for now
		_ = json.Unmarshal(metaJSON, &metadata)
	}
	conditionsReq := params.GetJSON("conditions")
	var conditions *map[string]interface{}
	if len(conditionsReq) > 0 {
		// marshal the conditions into the Map
		conditionsJSON, _ := json.Marshal(conditionsReq) // nolint: errchkjson // ignore for now
		_ = json.Unmarshal(conditionsJSON, &conditions)
	}

	// Record a new transaction (get the hex from parameters)a
	var err error
	var transactions []*bux.Transaction
	if transactions, err = a.Services.Bux.GetTransactions(
		req.Context(),
		reqXPub,
		metadata,
		conditions,
	); err != nil {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, err.Error())
		return
	}

	// Return response
	apirouter.ReturnResponse(w, req, http.StatusOK, bux.DisplayModels(transactions))
}
