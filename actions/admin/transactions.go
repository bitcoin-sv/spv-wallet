package admin

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/actions/common"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/internal/query"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
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

	contract := mappings.MapToTransactionContractForAdmin(transaction)
	c.JSON(http.StatusOK, contract)
}

// adminSearchTxs will fetch a list of transactions filtered by metadata
// Search for transactions filtering by metadata godoc
// @Summary		Search for transactions
// @Description	Fetches a list of transactions filtered by metadata and other criteria
// @Tags		Admin
// @Produce		json
// @Param		SwaggerCommonParams query swagger.CommonFilteringQueryParams false "Supports options for pagination and sorting to streamline data exploration and analysis"
// @Param		AdminTransactionFilter query filter.AdminTransactionFilter false "Supports targeted resource searches with filters"
// @Param		id query string false "Transaction ID"
// @Param		hex query string false "Transaction hex"
// @Param		blockHash query string false "Hash of the block containing the transaction"
// @Param		blockHeight query integer false "Height of the block containing the transaction"
// @Param		fee query integer false "Transaction fee"
// @Param		numberOfInputs query integer false "Number of inputs in the transaction"
// @Param		numberOfOutputs query integer false "Number of outputs in the transaction"
// @Param		draftId query string false "Draft ID associated with the transaction"
// @Param		totalValue query integer false "Total value of the transaction in satoshis"
// @Param		status query string false "Status of the transaction (e.g., 'confirmed', 'pending')"
// @Param		xpubId query string false "XPub ID associated with the transaction"
// @Success		200 {object} response.PageModel[response.Transaction] "List of transactions with pagination details"
// @Failure		400 "Bad request - Invalid query parameters"
// @Failure		500 "Internal server error - Error while searching for transactions"
// @Router		/api/v1/admin/transactions [get]
// @Security	x-auth-xpub
func adminSearchTxs(c *gin.Context, _ *reqctx.AdminContext) {
	logger := reqctx.Logger(c)

	searchParams, err := query.ParseSearchParams[filter.AdminTransactionFilter](c)
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotParseQueryParams.WithTrace(err), logger)
		return
	}

	queryParams := prepareQueryParams(c, searchParams)

	transactions, err := fetchTransactions(c, queryParams)
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrFetchTransactions.Wrap(err), logger)
		return
	}

	count, err := countTransactions(c, queryParams)
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCouldNotCountTransactions.WithTrace(err), logger)
		return
	}

	transactionContracts := common.MapToTypeContracts(transactions, mappings.MapToTransactionContractForAdmin)

	result := response.PageModel[response.Transaction]{
		Content: transactionContracts,
		Page:    common.GetPageDescriptionFromSearchParams(queryParams.PageOptions, count),
	}

	c.JSON(http.StatusOK, result)
}
