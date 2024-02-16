package accesskeys

import (
	"net/http"

	"github.com/bitcoin-sv/bux"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// get will get an existing model
// Get access key godoc
// @Summary		Get access key
// @Description	Get access key
// @Tags		Access-key
// @Produce		json
// @Param		id query string true "id"
// @Success		200
// @Router		/v1/access-key [get]
// @Security	spv-wallet-auth-xpub
func (a *Action) get(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	reqXPubID, _ := bux.GetXpubIDFromRequest(req)

	// Parse the params
	params := apirouter.GetParams(req)
	id := params.GetString("id")

	if id == "" {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, bux.ErrMissingFieldID)
		return
	}

	// Get access key
	accessKey, err := a.Services.SPV.GetAccessKey(
		req.Context(), reqXPubID, id,
	)
	if err != nil {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, err.Error())
		return
	}

	if accessKey.XpubID != reqXPubID {
		apirouter.ReturnResponse(w, req, http.StatusForbidden, "unauthorized")
		return
	}

	contract := mappings.MapToAccessKeyContract(accessKey)

	// Return response
	apirouter.ReturnResponse(w, req, http.StatusOK, contract)
}
