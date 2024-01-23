package destinations

import (
	"net/http"

	"github.com/BuxOrg/bux"
	"github.com/BuxOrg/bux-server/actions"
	"github.com/BuxOrg/bux-server/mappings"
	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// update will update an existing model
// Update Destination godoc
// @Summary		Update destination
// @Description	Update destination
// @Tags		Destinations
// @Produce		json
// @Param		id path string false "Destination ID"
// @Param		address path string false "Destination Address"
// @Param		locking_script path string false "Destination Locking Script"
// @Param		metadata body string true "Destination Metadata"
// @Success		200
// @Router		/v1/destination [patch]
// @Security	bux-auth-xpub
func (a *Action) update(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	reqXPubID, _ := bux.GetXpubIDFromRequest(req)

	// Parse the params
	params := apirouter.GetParams(req)
	id := params.GetString("id")
	address := params.GetString("address")
	lockingScript := params.GetString("locking_script")
	if id == "" && address == "" && lockingScript == "" {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, "One of the fields is required: id, address or lockingScript")
		return
	}
	metadata := params.GetJSON(actions.MetadataField)

	// Get the destination
	var destination *bux.Destination
	var err error
	if id != "" {
		destination, err = a.Services.Bux.UpdateDestinationMetadataByID(
			req.Context(), reqXPubID, id, metadata,
		)
	} else if address != "" {
		destination, err = a.Services.Bux.UpdateDestinationMetadataByAddress(
			req.Context(), reqXPubID, address, metadata,
		)
	} else {
		destination, err = a.Services.Bux.UpdateDestinationMetadataByLockingScript(
			req.Context(), reqXPubID, lockingScript, metadata,
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
