package accesskeys

import (
	"net/http"

	"github.com/BuxOrg/bux"
	"github.com/BuxOrg/bux-server/mappings"
	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// revoke will revoke the intended model by id
// Revoke access key godoc
// @Summary		Revoke access key
// @Description	Revoke access key
// @Tags		Access-key
// @Produce		json
// @Param		id query string true "id"
// @Success		201
// @Router		/v1/access-key [delete]
// @Security	bux-auth-xpub
func (a *Action) revoke(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	reqXPub, _ := bux.GetXpubFromRequest(req)

	// Parse the params
	params := apirouter.GetParams(req)
	id := params.GetString("id")

	if id == "" {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, bux.ErrMissingFieldID)
		return
	}

	// Create a new accessKey
	accessKey, err := a.Services.Bux.RevokeAccessKey(
		req.Context(),
		reqXPub,
		id,
	)
	if err != nil {
		apirouter.ReturnResponse(w, req, http.StatusUnprocessableEntity, err.Error())
		return
	}

	contract := mappings.MapToAccessKeyContract(accessKey)

	// Return response
	apirouter.ReturnResponse(w, req, http.StatusCreated, bux.DisplayModels(contract))
}
