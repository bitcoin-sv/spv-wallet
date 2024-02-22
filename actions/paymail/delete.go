package pmail

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// delete will remove the intended model
// Delete Paymail godoc
// @Summary		Delete paymail
// @Description	Delete paymail
// @Tags		Paymails
// @Param		address query string true "address"
// @Produce		json
// @Success		200
// @Router		/v1/paymail [delete]
// @Security	x-auth-xpub
func (a *Action) delete(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	// Parse the params
	params := apirouter.GetParams(req)

	// params
	address := params.GetString("address") // the full paymail address

	opts := a.Services.SpvWalletEngine.DefaultModelOptions()

	// Delete a new paymail address
	err := a.Services.SpvWalletEngine.DeletePaymailAddress(
		req.Context(), address, opts...,
	)
	if err != nil {
		apirouter.ReturnResponse(w, req, http.StatusUnprocessableEntity, err.Error())
		return
	}

	// Return response
	apirouter.ReturnResponse(w, req, http.StatusOK, nil)
}
