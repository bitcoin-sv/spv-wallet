package admin

import (
	"net/http"

	"github.com/BuxOrg/bux"
	"github.com/BuxOrg/bux-server/actions"
	"github.com/julienschmidt/httprouter"
	apirouter "github.com/mrz1836/go-api-router"
)

// transactionsSearch will fetch a list of transactions filtered by metadata
// Search for transactions filtering by metadata godoc
// @Summary      Search for transactions
// @Description  Search for transactions
// @Tags		 Admin
// @Produce      json
// @Param        page query int false "page"
// @Param        page_size query int false "page_size"
// @Param        order_by_field query string false "order_by_field"
// @Param        sort_direction query string false "sort_direction"
// @Param 		 metadata query string false "Metadata filter"
// @Param 		 conditions query string false "Conditions filter"
// @Success      200
// @Router       /v1/admin/transactions/search [post]
// @Security bux-auth-xpub
func (a *Action) transactionsSearch(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	// Parse the params
	params := apirouter.GetParams(req)
	queryParams, metadata, conditions, err := actions.GetQueryParameters(params)
	if err != nil {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, err.Error())
		return
	}

	var transactions []*bux.Transaction
	if transactions, err = a.Services.Bux.GetTransactions(
		req.Context(),
		metadata,
		conditions,
		queryParams,
	); err != nil {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, err.Error())
		return
	}

	// Return response
	apirouter.ReturnResponse(w, req, http.StatusOK, transactions)
}

// transactionsCount will count all transactions filtered by metadata
// Count transactions filtering by metadata godoc
// @Summary      Count transactions
// @Description  Count transactions
// @Tags		 Admin
// @Produce      json
// @Param metadata query string false "Metadata filter"
// @Param conditions query string false "Conditions filter"
// @Success      200
// @Router       /v1/admin/transactions/count [post]
// @Security bux-auth-xpub
func (a *Action) transactionsCount(w http.ResponseWriter, req *http.Request, _ httprouter.Params) {
	// Parse the params
	params := apirouter.GetParams(req)
	_, metadata, conditions, err := actions.GetQueryParameters(params)
	if err != nil {
		apirouter.ReturnResponse(w, req, http.StatusExpectationFailed, err.Error())
		return
	}

	var count int64
	if count, err = a.Services.Bux.GetTransactionsCount(
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
