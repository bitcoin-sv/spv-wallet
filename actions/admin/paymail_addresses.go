package admin

import (
	"net/http"

	"github.com/BuxOrg/bux"
	"github.com/BuxOrg/bux-server/actions"
	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// paymailAddressesSearch will fetch a list of paymail addresses filtered by metadata
func (a *Action) paymailAddressesSearch(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	// Parse the params
	params := apirouter.GetParams(req)
	queryParams, metadata, conditions, err := actions.GetQueryParameters(params)
	if err != nil {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, err.Error())
		return
	}

	var paymailAddresses []*bux.PaymailAddress
	if paymailAddresses, err = a.Services.Bux.GetPaymailAddresses(
		req.Context(),
		metadata,
		conditions,
		queryParams,
	); err != nil {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, err.Error())
		return
	}

	// Return response
	apirouter.ReturnResponse(w, req, http.StatusOK, paymailAddresses)
}

// paymailAddressesCount will count all paymail addresses filtered by metadata
func (a *Action) paymailAddressesCount(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	// Parse the params
	params := apirouter.GetParams(req)
	_, metadata, conditions, err := actions.GetQueryParameters(params)
	if err != nil {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, err.Error())
		return
	}

	var count int64
	if count, err = a.Services.Bux.GetPaymailAddressesCount(
		req.Context(),
		metadata,
		conditions,
	); err != nil {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, err.Error())
		return
	}

	// Return response
	apirouter.ReturnResponse(w, req, http.StatusOK, count)
}

// paymailCreateAddress will create a new paymail address
func (a *Action) paymailCreateAddress(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	// Parse the params
	params := apirouter.GetParams(req)

	xpub := params.GetString("xpub")
	address := params.GetString("address")
	publicName := params.GetString("public_name")
	avatar := params.GetString("avatar")

	if xpub == "" {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, "xpub is required")
		return
	}
	if address == "" {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, "address is required")
		return
	}

	_, metadata, _, err := actions.GetQueryParameters(params)
	if err != nil {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, err.Error())
		return
	}

	opts := a.Services.Bux.DefaultModelOptions()

	if metadata != nil {
		opts = append(opts, bux.WithMetadatas(*metadata))
	}

	var paymailAddress *bux.PaymailAddress
	paymailAddress, err = a.Services.Bux.NewPaymailAddress(req.Context(), xpub, address, publicName, avatar, opts...)
	if err != nil {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, err.Error())
		return
	}

	// Return response
	apirouter.ReturnResponse(w, req, http.StatusOK, paymailAddress)
}

// paymailDeleteAddress will delete a paymail address
func (a *Action) paymailDeleteAddress(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	// Parse the params
	params := apirouter.GetParams(req)
	address := params.GetString("address")

	if address == "" {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, "address is required")
		return
	}

	opts := a.Services.Bux.DefaultModelOptions()

	// Delete a new paymail address
	err := a.Services.Bux.DeletePaymailAddress(req.Context(), address, opts...)
	if err != nil {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, err.Error())
		return
	}

	// Return response
	apirouter.ReturnResponse(w, req, http.StatusOK, true)
}
