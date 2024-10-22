package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bitcoin-sv/spv-wallet/actions/common"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/internal/query"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
)

// adminGetTxByID fetches a transaction by id for admins
// @Summary		Get transaction by id for admins
// @Description	Get transaction by id for admins
// @Tags		Admin
// @Produce		json
// @Param		id path string true "Transaction ID"
// @Success		200 {object} response.Transaction "Transaction"
// @Failure		400	"Bad request - Transaction not found or error in data fetching"
// @Failure 	500	"Internal Server Error - Error while fetching transaction"
// @Router		/api/v1/admin/transactions/{id} [get]
// @Security	x-auth-xpub
func adminGetTxByID(c *gin.Context, _ *reqctx.AdminContext) {
	logger := reqctx.Logger(c)
	id := c.Param("id")

	transaction, err := reqctx.Engine(c).GetAdminTransaction(c, id)
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCouldNotFindTransaction.WithTrace(err), logger)
		return
	}
	if transaction == nil {
		spverrors.ErrorResponse(c, spverrors.ErrCouldNotFindTransaction, logger)
		return
	}

	contract := mappings.MapToTransactionContract(transaction)
	c.JSON(http.StatusOK, contract)
}

// adminSearchTxs will fetch a list of transactions filtered by metadata
// Search for transactions filtering by metadata godoc
// @Summary		Search for transactions
// @Description	Search for transactions
// @Tags		Admin
// @Produce		json
// @Param		metadata query string false "Filter by metadata in the form of key-value pairs"
// @Param		conditions query string false "Additional conditions for filtering, in URL-encoded JSON"
// @Param		queryParams query string false "Pagination and sorting options"
// @Success		200 {object} []response.Transaction "List of transactions"
// @Failure		400 "Bad request - Error while parsing query parameters"
// @Failure 	500 "Internal server error - Error while searching for transactions"
// @Router		/api/v1/admin/transactions [get]
// @Security	x-auth-xpub
func adminSearchTxs(c *gin.Context, _ *reqctx.AdminContext) {
	logger := reqctx.Logger(c)

	searchParams, err := query.ParseSearchParams[filter.TransactionFilter](c)
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotParseQueryParams.WithTrace(err), logger)
		return
	}

	conditions := searchParams.Conditions.ToDbConditions()
	metadata := mappings.MapToMetadata(searchParams.Metadata)
	pageOptions := mappings.MapToDbQueryParams(&searchParams.Page)

	transactions, err := reqctx.Engine(c).GetTransactions(
		c.Request.Context(),
		metadata,
		conditions,
		pageOptions,
	)
	if err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	}

	count, err := reqctx.Engine(c).GetTransactionsCount(
		c.Request.Context(),
		metadata,
		conditions,
	)
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCouldNotCountTransactions.WithTrace(err), logger)
		return
	}

	contracts := common.MapToTypeContracts(transactions, mappings.MapToTransactionContractForAdmin)
	result := response.PageModel[response.Transaction]{
		Content: contracts,
		Page:    common.GetPageDescriptionFromSearchParams(pageOptions, count),
	}

	c.JSON(http.StatusOK, result)
}
