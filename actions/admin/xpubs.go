package admin

import (
	"net/http"

	"github.com/BuxOrg/bux"
	"github.com/BuxOrg/bux-server/actions"
	"github.com/BuxOrg/bux-server/mappings"
	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// xpubsSearch will fetch a list of xpubs filtered by metadata
// Search for xpubs filtering by metadata godoc
// @Summary		Search for xpubs
// @Description	Search for xpubs
// @Tags		Admin
// @Produce		json
// @Param		page query int false "page"
// @Param		page_size query int false "page_size"
// @Param		order_by_field query string false "order_by_field"
// @Param		sort_direction query string false "sort_direction"
// @Param		metadata query string false "Metadata filter"
// @Param		conditions query string false "Conditions filter"
// @Success		200
// @Router		/v1/admin/xpubs/search [post]
// @Security	bux-auth-xpub
func (a *Action) xpubsSearch(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	// Parse the params
	params := apirouter.GetParams(req)
	queryParams, metadataModel, conditions, err := actions.GetQueryParameters(params)
	metadata := mappings.MapToBuxMetadata(metadataModel)

	if err != nil {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, err.Error())
		return
	}

	var xpubs []*bux.Xpub
	if xpubs, err = a.Services.Bux.GetXPubs(
		req.Context(),
		metadata,
		conditions,
		queryParams,
	); err != nil {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, err.Error())
		return
	}

	// Return response
	apirouter.ReturnResponse(w, req, http.StatusOK, xpubs)
}

// xpubsCount will count all xpubs filtered by metadata
// Count xpubs filtering by metadata godoc
// @Summary		Count xpubs
// @Description	Count xpubs
// @Tags		Admin
// @Produce		json
// @Param		metadata query string false "Metadata filter"
// @Param		conditions query string false "Conditions filter"
// @Success		200
// @Router		/v1/admin/xpubs/count [post]
// @Security	bux-auth-xpub
func (a *Action) xpubsCount(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	// Parse the params
	params := apirouter.GetParams(req)
	_, metadataModel, conditions, err := actions.GetQueryParameters(params)
	metadata := mappings.MapToBuxMetadata(metadataModel)
	if err != nil {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, err.Error())
		return
	}

	var count int64
	if count, err = a.Services.Bux.GetXPubsCount(
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
