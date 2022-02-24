package pmail

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// create will make a new model using the services defined in the action object
func (a *Action) delete(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	// Parse the params
	params := apirouter.GetParams(req)

	// params
	address := params.GetString("address") // the full paymail address

	opts := a.Services.Bux.DefaultModelOptions()

	// Create a new paymail address
	err := deletePaymailAddress(
		req.Context(), address, opts...,
	)
	if err != nil {
		apirouter.ReturnResponse(w, req, http.StatusUnprocessableEntity, err.Error())
		return
	}

	// Return response
	apirouter.ReturnResponse(w, req, http.StatusCreated, nil)
}
