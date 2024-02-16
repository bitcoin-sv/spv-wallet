package admin

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/actions"
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// destinationsSearch will fetch a list of destinations filtered by metadata
// Search for destinations filtering by metadata godoc
// @Summary		Search for destinations
// @Description	Search for destinations
// @Tags		Admin
// @Produce		json
// @Param		page query int false "page"
// @Param		page_size query int false "page_size"
// @Param		order_by_field query string false "order_by_field"
// @Param		sort_direction query string false "sort_direction"
// @Param		metadata query string false "Metadata filter"
// @Param		conditions query string false "Conditions filter"
// @Success		200
// @Router		/v1/admin/destinations/search [post]
// @Security	x-auth-xpub
func (a *Action) destinationsSearch(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	// Parse the params
	params := apirouter.GetParams(req)
	queryParams, metadataModel, conditions, err := actions.GetQueryParameters(params)
	metadata := mappings.MapToSPVMetadata(metadataModel)
	if err != nil {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, err.Error())
		return
	}

	var destinations []*engine.Destination
	if destinations, err = a.Services.SPV.GetDestinations(
		req.Context(),
		metadata,
		conditions,
		queryParams,
	); err != nil {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, err.Error())
		return
	}

	// Return response
	apirouter.ReturnResponse(w, req, http.StatusOK, destinations)
}

// destinationsCount will count all destinations filtered by metadata
// Count destinations filtering by metadata godoc
// @Summary		Count destinations
// @Description	Count destinations
// @Tags		Admin
// @Produce		json
// @Param		metadata query string false "Metadata filter"
// @Param		conditions query string false "Conditions filter"
// @Success		200
// @Router		/v1/admin/destinations/count [post]
// @Security	x-auth-xpub
func (a *Action) destinationsCount(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	// Parse the params
	params := apirouter.GetParams(req)
	_, metadataModel, conditions, err := actions.GetQueryParameters(params)
	metadata := mappings.MapToSPVMetadata(metadataModel)
	if err != nil {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, err.Error())
		return
	}

	var count int64
	if count, err = a.Services.SPV.GetDestinationsCount(
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
