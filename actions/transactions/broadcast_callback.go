package transactions

import (
	"encoding/json"
	"net/http"

	"github.com/bitcoin-sv/go-broadcast-client/broadcast"
	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// broadcastCallback will handle a broadcastCallback call from the broadcast api
// Broadcast Callback godoc
// @Summary		Broadcast Callback
// @Tags		Transactions
// @Param 		transaction body broadcast.SubmittedTx true "transaction"
// @Success		200
// @Router		/transaction/broadcast/callback [post]
// @Security	callback-auth
func (a *Action) broadcastCallback(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	var resp *broadcast.SubmittedTx

	err := json.NewDecoder(req.Body).Decode(&resp)
	if err != nil {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, err.Error())
		return
	}

	defer func() {
		if err = req.Body.Close(); err != nil {
			a.Services.Logger.Err(err).Msg("failed to close request body")
		}
	}()

	err = a.Services.SPV.UpdateTransaction(req.Context(), resp)
	if err != nil {
		a.Services.Logger.Err(err).Msgf("failed to update transaction - tx: %v", resp)
		apirouter.ReturnResponse(w, req, http.StatusInternalServerError, "")
		return
	}

	// Return response
	apirouter.ReturnResponse(w, req, http.StatusOK, "")
}
