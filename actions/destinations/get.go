package destinations

import (
	"net/http"

	"github.com/BuxOrg/bux"
	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// get will get an existing model
func (a *Action) get(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {

	reqXPub, _ := bux.GetXpubFromRequest(req)

	// Parse the params
	params := apirouter.GetParams(req)
	id := params.GetString("id")
	address := params.GetString("address")
	lockingScript := params.GetString("locking_script")
	if id == "" && address == "" && lockingScript == "" {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, bux.ErrMissingFieldID)
		return
	}

	// Get the destination
	var destination *bux.Destination
	var err error
	if id != "" {
		destination, err = a.Services.Bux.GetDestinationByID(
			req.Context(), reqXPub, id,
		)
	} else if address != "" {
		destination, err = a.Services.Bux.GetDestinationByAddress(
			req.Context(), reqXPub, address,
		)
	} else {
		destination, err = a.Services.Bux.GetDestinationByLockingScript(
			req.Context(), reqXPub, lockingScript,
		)
	}
	if err != nil {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, err.Error())
		return
	}

	// Return response
	apirouter.ReturnResponse(w, req, http.StatusOK, bux.DisplayModels(destination))
}
