package destinations

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// get will get an existing model
// Get Destination godoc
// @Summary		Get a destination
// @Description	Get a destination
// @Tags		Destinations
// @Produce		json
// @Param		id query string false "Destination ID"
// @Param		address query string false "Destination address"
// @Param		locking_script query string false "Destination locking script"
// @Success		200
// @Router		/v1/destination [get]
// @Security	x-auth-xpub
func (a *Action) get(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	reqXPubID, _ := engine.GetXpubIDFromRequest(req)

	// Parse the params
	params := apirouter.GetParams(req)
	id := params.GetString("id")
	address := params.GetString("address")
	lockingScript := params.GetString("locking_script")
	if id == "" && address == "" && lockingScript == "" {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, engine.ErrMissingFieldID)
		return
	}

	// Get the destination
	var destination *engine.Destination
	var err error
	if id != "" {
		destination, err = a.Services.SpvWalletEngine.GetDestinationByID(
			req.Context(), reqXPubID, id,
		)
	} else if address != "" {
		destination, err = a.Services.SpvWalletEngine.GetDestinationByAddress(
			req.Context(), reqXPubID, address,
		)
	} else {
		destination, err = a.Services.SpvWalletEngine.GetDestinationByLockingScript(
			req.Context(), reqXPubID, lockingScript,
		)
	}
	if err != nil {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, err.Error())
		return
	}

	contract := mappings.MapToDestinationContract(destination)

	// Return response
	apirouter.ReturnResponse(w, req, http.StatusOK, contract)
}
