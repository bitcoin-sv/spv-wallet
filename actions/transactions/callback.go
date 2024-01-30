package transactions

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// callback will handle a callback call from the broadcast api
// Broadcast Callback godoc
// @Summary		Broadcast Callback
// @Tags		Transactions
// @Param 		transaction body broadcast.SubmittedTx true "transaction"
// @Success		200
// @Router		/v1/transaction/callback [post]
// @Security	callback-auth
func (a *Action) callback(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var transaction *broadcast.SubmittedTx

	err := json.NewDecoder(req.Body).Decode(&transaction)
	if err != nil {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, err.Error())
		return
	}

	defer req.Body.Close()

	log.Printf("Received Tx: %+v", transaction)

	// Return response
	apirouter.ReturnResponse(w, req, http.StatusOK, "")
}
