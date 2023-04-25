package admin

import (
	"net/http"

	"github.com/BuxOrg/bux"
	"github.com/BuxOrg/bux-server/actions"
	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// utxosSearch will fetch a list of utxos filtered by metadata
// Search for utxos filtering by metadata godoc
// @Summary      Search for utxos
// @Description  Search for utxos
// @Tags		 Admin
// @Produce      json
// @Param       	page query int false "page"
// @Param       	page_size query int false "page_size"
// @Param       	order_by_field query string false "order_by_field"
// @Param       	sort_direction query string false "sort_direction"
// @Param metadata query string false "Metadata filter"
// @Param conditions query string false "Conditions filter"
// @Success      200
// @Router       /v1/admin/utxos/search [post]
// @Security bux-auth-xpub
func (a *Action) utxosSearch(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	// Parse the params
	params := apirouter.GetParams(req)
	queryParams, metadata, conditions, err := actions.GetQueryParameters(params)
	if err != nil {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, err.Error())
		return
	}

	var utxos []*bux.Utxo
	if utxos, err = a.Services.Bux.GetUtxos(
		req.Context(),
		metadata,
		conditions,
		queryParams,
	); err != nil {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, err.Error())
		return
	}

	// Return response
	apirouter.ReturnResponse(w, req, http.StatusOK, utxos)
}

// utxosCount will count all utxos filtered by metadata
// Count utxos filtering by metadata godoc
// @Summary      Count utxos
// @Description  Count utxos
// @Tags		 Admin
// @Produce      json
// @Param metadata query string false "Metadata filter"
// @Param conditions query string false "Conditions filter"
// @Success      200
// @Router       /v1/admin/utxos/count [post]
// @Security bux-auth-xpub
func (a *Action) utxosCount(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	// Parse the params
	params := apirouter.GetParams(req)
	_, metadata, conditions, err := actions.GetQueryParameters(params)
	if err != nil {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, err.Error())
		return
	}

	var count int64
	if count, err = a.Services.Bux.GetUtxosCount(
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
