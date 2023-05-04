package destinations

import (
	"net/http"

	"github.com/BuxOrg/bux"
	"github.com/BuxOrg/bux-server/actions"
	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// count will fetch a count of destinations filtered by metadata
// Count Destinations godoc
// @Summary		Count Destinations
// @Description	Count Destinations
// @Tags		Destinations
// @Param		metadata query string false "metadata"
// @Param		condition query string false "condition"
// @Produce		json
// @Success		200
// @Router		/v1/destination/count [post]
// @Security	bux-auth-xpub
func (a *Action) count(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	reqXPubID, _ := bux.GetXpubIDFromRequest(req)

	// Parse the params
	params := apirouter.GetParams(req)
	_, metadata, conditions, err := actions.GetQueryParameters(params)
	if err != nil {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, err.Error())
		return
	}

	// Record a new transaction (get the hex from parameters)
	var count int64
	if count, err = a.Services.Bux.GetDestinationsByXpubIDCount(
		req.Context(),
		reqXPubID,
		metadata,
		conditions,
	); err != nil {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, err.Error())
		return
	}

	// Return response
	apirouter.ReturnResponse(w, req, http.StatusOK, count)
}
