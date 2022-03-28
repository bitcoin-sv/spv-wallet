package transactions

import (
	"net/http"

	"github.com/BuxOrg/bux"
	"github.com/BuxOrg/bux-server/actions"
	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// get will fetch a transaction
func (a *Action) update(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {

	// Get the xPub from the request (via authentication)
	reqXPubID, _ := bux.GetXpubIDFromRequest(req)

	// Parse the params
	params := apirouter.GetParams(req)
	metadata := params.GetJSON(actions.MetadataField)

	// Get a transaction by ID
	transaction, err := a.Services.Bux.UpdateTransactionMetadata(
		req.Context(),
		reqXPubID,
		params.GetString("id"),
		metadata,
	)
	if err != nil {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, err.Error())
		return
	} else if transaction == nil {
		apirouter.ReturnResponse(w, req, http.StatusNotFound, "")
	} else if !transaction.IsXpubIDAssociated(reqXPubID) {
		apirouter.ReturnResponse(w, req, http.StatusForbidden, "unauthorized")
		return
	}

	// Return response
	apirouter.ReturnResponse(w, req, http.StatusOK, bux.DisplayModels(transaction))
}
