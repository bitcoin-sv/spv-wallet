package utxos

import (
	"net/http"

	"github.com/BuxOrg/bux"
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
// @Security	spv-wallet-auth-xpub
func (a *Action) unreserve(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	reqXPubID, _ := bux.GetXpubIDFromRequest(req)
	params := apirouter.GetParams(req)

	err := a.Services.SPV.UnReserveUtxos(
		req.Context(),
		reqXPubID,
		params.GetString(bux.ReferenceIDField),
	)
	if err != nil {
		apirouter.ReturnResponse(w, req, http.StatusUnprocessableEntity, err.Error())
		return
	}

	apirouter.ReturnResponse(w, req, http.StatusNoContent, 0)
}
