package xpubs

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// get will get an existing model
// Get xPub godoc
// @Summary		Get xPub
// @Description	Get xPub
// @Tags		xPub
// @Produce		json
// @Param		key query string false "key"
// @Success		200
// @Router		/v1/xpub [get]
// @Security	x-auth-xpub
func (a *Action) get(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	reqXPub, _ := engine.GetXpubFromRequest(req)
	reqXPubID, _ := engine.GetXpubIDFromRequest(req)

	// Parse the params
	params := apirouter.GetParams(req)
	key := params.GetString("key")
	if key != "" {
		if isAdmin, ok := engine.IsAdminRequest(req); !isAdmin || !ok {
			apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, engine.ErrNotAdminKey)
			return
		}
	} else {
		key = reqXPub
	}

	// Get an xPub
	var xPub *engine.Xpub
	var err error
	if key != "" {
		xPub, err = a.Services.SPV.GetXpub(
			req.Context(), key,
		)
	} else {
		xPub, err = a.Services.SPV.GetXpubByID(
			req.Context(), reqXPubID,
		)
	}
	if err != nil {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, err.Error())
		return
	}

	signed := req.Context().Value("auth_signed")
	if signed == nil || !signed.(bool) || reqXPub == "" {
		xPub.RemovePrivateData()
	}

	contract := mappings.MapToXpubContract(xPub)

	// Return response
	apirouter.ReturnResponse(w, req, http.StatusOK, contract)
}
