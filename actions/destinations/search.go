package destinations

import (
	"net/http"

	"github.com/BuxOrg/bux"
	buxmodels "github.com/BuxOrg/bux-models"
	"github.com/BuxOrg/bux-server/actions"
	"github.com/BuxOrg/bux-server/mappings"
	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// search will fetch a list of destinations filtered by metadata
// Search Destination godoc
// @Summary		Search for a destination
// @Description	Search for a destination
// @Tags		Destinations
// @Produce		json
// @Param		page query int false "page"
// @Param		page_size query int false "page_size"
// @Param		order_by_field query string false "order_by_field"
// @Param		sort_direction query string false "sort_direction"
// @Param		metadata query string false "metadata"
// @Param		condition query string false "condition"
// @Success		200
// @Router		/v1/destination/search [get]
// @Security	bux-auth-xpub
func (a *Action) search(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	reqXPubID, _ := bux.GetXpubIDFromRequest(req)

	// Parse the params
	params := apirouter.GetParams(req)
	queryParams, metadataModel, conditions, err := actions.GetQueryParameters(params)
	metadata := mappings.MapToBuxMetadata(metadataModel)
	if err != nil {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, err.Error())
		return
	}

	// Record a new transaction (get the hex from parameters)a
	var destinations []*bux.Destination
	if destinations, err = a.Services.Bux.GetDestinationsByXpubID(
		req.Context(),
		reqXPubID,
		metadata,
		conditions,
		queryParams,
	); err != nil {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, err.Error())
		return
	}

	contracts := make([]*buxmodels.Destination, 0)
	for _, destination := range destinations {
		contracts = append(contracts, mappings.MapToDestinationContract(destination))
	}

	// Return response
	apirouter.ReturnResponse(w, req, http.StatusOK, contracts)
}
