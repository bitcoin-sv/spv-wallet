package utxos

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine"

	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// unreserve remove the reservation on the utxos for the given draft ID
// Unreserve UTXOs godoc
// @Summary		Unreserve UTXOs
// @Description	Unreserve UTXOs
// @Tags		UTXO
// @Param		reference_id query string false "draft tx id"
// @Success		201
// @Router		/v1/utxo/unreserve [patch]
// @Security	x-auth-xpub
func (a *Action) unreserve(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	reqXPubID, _ := engine.GetXpubIDFromRequest(req)
	params := apirouter.GetParams(req)

	err := a.Services.SpvWalletEngine.UnReserveUtxos(
		req.Context(),
		reqXPubID,
		params.GetString(engine.ReferenceIDField),
	)
	if err != nil {
		apirouter.ReturnResponse(w, req, http.StatusUnprocessableEntity, err.Error())
		return
	}

	apirouter.ReturnResponse(w, req, http.StatusNoContent, 0)
}
