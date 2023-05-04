package admin

import (
	"net/http"

	"github.com/BuxOrg/bux"
	"github.com/BuxOrg/bux-server/actions"
	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// blockHeadersSearch will fetch a list of block headers filtered by metadata
// Search for block headers godoc
// @Summary		Search for block headers
// @Description	Search for block headers
// @Tags		Admin
// @Param		page query int false "page"
// @Param		page_size query int false "page_size"
// @Param		order_by_field query string false "order_by_field"
// @Param		sort_direction query string false "sort_direction"
// @Param		metadata query string false "Metadata filter"
// @Param		conditions query string false "Conditions filter"
// @Produce		json
// @Success		200
// @Router		/v1/admin/block-headers/search [post]
// @Security	bux-auth-xpub
func (a *Action) blockHeadersSearch(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	// Parse the params
	params := apirouter.GetParams(req)
	queryParams, metadata, conditions, err := actions.GetQueryParameters(params)
	if err != nil {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, err.Error())
		return
	}

	var blockHeaders []*bux.BlockHeader
	if blockHeaders, err = a.Services.Bux.GetBlockHeaders(
		req.Context(),
		metadata,
		conditions,
		queryParams,
	); err != nil {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, err.Error())
		return
	}

	// Return response
	apirouter.ReturnResponse(w, req, http.StatusOK, blockHeaders)
}

// blockHeadersCount will count all block headers filtered by metadata
// Get block headers count headers godoc
// @Summary		Get block headers count
// @Description	Get block headers count
// @Tags		Admin
// @Param		metadata query string false "Metadata filter"
// @Param		conditions query string false "Conditions filter"
// @Produce		json
// @Success		200
// @Router		/v1/admin/block-headers/count [post]
// @Security	bux-auth-xpub
func (a *Action) blockHeadersCount(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	// Parse the params
	params := apirouter.GetParams(req)
	_, metadata, conditions, err := actions.GetQueryParameters(params)
	if err != nil {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, err.Error())
		return
	}

	var count int64
	if count, err = a.Services.Bux.GetBlockHeadersCount(
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
