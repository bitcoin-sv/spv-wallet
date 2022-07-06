package utxos

import (
	"net/http"

	"github.com/BuxOrg/bux"
	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// get will fetch a given utxo according to conditions
func (a *Action) get(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	reqXPubID, _ := bux.GetXpubIDFromRequest(req)

	// Parse the params
	params := apirouter.GetParams(req)
	txID := params.GetString("tx_id")
	outputIndex := uint32(params.GetUint64("output_index"))

	// Get a utxo using a xPub
	utxo, err := a.Services.Bux.GetUtxo(
		req.Context(),
		reqXPubID,
		txID,
		outputIndex,
	)
	if err != nil {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, err.Error())
		return
	}

	// Return response
	apirouter.ReturnResponse(w, req, http.StatusOK, bux.DisplayModels(utxo))
}
