package admin

import (
	"net/http"

	"github.com/BuxOrg/bux"
	"github.com/BuxOrg/bux-server/actions"
	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// accessKeysSearch will fetch a list of access keys filtered by metadata
// Access Keys Search godoc
// @Summary      Access Keys Search
// @Description  Access Keys Search
// @Tags		 Admin
// @Produce      json
// @Param page query int false "page"
// @Param page_size query int false "page_size"
// @Param order_by_field query string false "order_by_field"
// @Param sort_direction query string false "sort_direction"
// @Param metadata query string false "Metadata filter"
// @Param conditions query string false "Conditions filter"
// @Success      200
// @Router       /v1/admin/access-keys/search [post]
// @Security bux-auth-xpub
func (a *Action) accessKeysSearch(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	// Parse the params
	params := apirouter.GetParams(req)
	queryParams, metadata, conditions, err := actions.GetQueryParameters(params)
	if err != nil {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, err.Error())
		return
	}

	var accessKeys []*bux.AccessKey
	if accessKeys, err = a.Services.Bux.GetAccessKeys(
		req.Context(),
		metadata,
		conditions,
		queryParams,
	); err != nil {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, err.Error())
		return
	}

	// Return response
	apirouter.ReturnResponse(w, req, http.StatusOK, accessKeys)
}

// accessKeysCount will count all access keys filtered by metadata
// Access Keys Count godoc
// @Summary      Access Keys Count
// @Description  Access Keys Count
// @Tags		 Admin
// @Produce      json
// @Param metadata query string false "Metadata filter"
// @Param conditions query string false "Conditions filter"
// @Success      200
// @Router       /v1/admin/access-keys/count [post]
// @Security bux-auth-xpub
func (a *Action) accessKeysCount(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	// Parse the params
	params := apirouter.GetParams(req)
	_, metadata, conditions, err := actions.GetQueryParameters(params)
	if err != nil {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, err.Error())
		return
	}

	var count int64
	if count, err = a.Services.Bux.GetAccessKeysCount(
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
