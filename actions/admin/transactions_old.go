package admin

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

// transactionsSearch will fetch a list of transactions filtered by metadata
// Search for transactions filtering by metadata godoc
// @Summary		Search for transactions
// @DeprecatedRouter	/v1/admin/transactions [post]
// @Description	Search for transactions
// @Tags		Admin
// @Produce		json
// @Param		SearchTransactions body filter.SearchTransactions false "Supports targeted resource searches with filters and metadata, plus options for pagination and sorting to streamline data exploration and analysis"
// @Success		200 {object} []models.Transaction "List of transactions"
// @Failure		400	"Bad request - Error while parsing SearchTransactions from request body"
// @Failure 	500	"Internal server error - Error while searching for transactions"
// @Router		/v1/admin/transactions/search [post]
// @Security	x-auth-xpub
func transactionsSearch(c *gin.Context, _ *reqctx.AdminContext) {
	logger := reqctx.Logger(c)
	var reqParams filter.SearchTransactions
	if err := c.Bind(&reqParams); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, logger)
		return
	}

	transactions, err := reqctx.Engine(c).GetTransactions(
		c.Request.Context(),
		mappings.MapToMetadata(reqParams.Metadata),
		reqParams.Conditions.ToDbConditions(),
		mappings.MapToQueryParams(reqParams.QueryParams),
	)
	if err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	}

	contracts := make([]*models.Transaction, 0)
	for _, transaction := range transactions {
		contracts = append(contracts, mappings.MapToOldTransactionContractForAdmin(transaction))
	}

	c.JSON(http.StatusOK, contracts)
}

// transactionsCount will count all transactions filtered by metadata
// Count transactions filtering by metadata godoc
// @DeprecatedRouter /v1/admin/transactions/count [post]
// @Summary		Count transactions
// @Description	Count transactions
// @Tags		Admin
// @Produce		json
// @Param		CountTransactions body filter.CountTransactions false "Enables filtering of elements to be counted"
// @Success		200	{number} int64 "Count of transactions"
// @Failure		400	"Bad request - Error while parsing CountTransactions from request body"
// @Failure 	500	"Internal Server Error - Error while fetching count of transactions"
// @Router		/v1/admin/transactions/count [post]
// @Security	x-auth-xpub
func transactionsCount(c *gin.Context, _ *reqctx.AdminContext) {
	logger := reqctx.Logger(c)
	var reqParams filter.CountTransactions
	if err := c.Bind(&reqParams); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, logger)
		return
	}

	count, err := reqctx.Engine(c).GetTransactionsCount(
		c.Request.Context(),
		mappings.MapToMetadata(reqParams.Metadata),
		reqParams.Conditions.ToDbConditions(),
	)
	if err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	}

	c.JSON(http.StatusOK, count)
}

// getTxAdminDeprecated fetches a transaction for admins using a deprecated path
// @DeprecatedRouter /v1/admin/transactions [get]
// @Summary		Get transaction by id for admins (Deprecated)
// @Description	This endpoint has been deprecated. Use /v1/admin/transactions/{id} instead.
// @Tags		Admin
// @Produce		json
// @Param		id query string true "id"
// @Success		200 {object} models.Transaction "Transaction"
// @Failure		400	"Bad request - Transaction not found or error in data fetching"
// @Failure 	500	"Internal Server Error - Error while fetching transaction"
// @Router		/v1/admin/transactions [get]
// @Security	x-auth-xpub
func getTxAdminByIDOld(c *gin.Context, _ *reqctx.AdminContext) {
	logger := reqctx.Logger(c)
	id := c.Param("id")

	transaction, err := reqctx.Engine(c).GetAdminTransaction(c.Request.Context(), id)
	if err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	}
	if transaction == nil {
		spverrors.ErrorResponse(c, spverrors.ErrCouldNotFindTransaction, logger)
		return
	}

	contract := mappings.MapToOldTransactionContract(transaction)
	c.JSON(http.StatusOK, contract)
}

// getTransactionsOld will fetch a list of transactions filtered by metadata
// Search for transactions filtering by metadata godoc
// @DeprecatedRouter /v1/admin/transactions [get]
// @Summary		Search for transactions
// @Description	Search for transactions
// @Tags		Admin
// @Produce		json
// @Param		metadata query string false "Filter by metadata in the form of key-value pairs"
// @Param		conditions query string false "Additional conditions for filtering, in URL-encoded JSON"
// @Param		queryParams query string false "Pagination and sorting options"
// @Success		200 {object} []models.Transaction "List of transactions"
// @Failure		400 "Bad request - Error while parsing query parameters"
// @Failure 	500 "Internal server error - Error while searching for transactions"
// @Router		/v1/admin/transactions [get]
// @Security	x-auth-xpub
func getTransactionsOld(c *gin.Context, _ *reqctx.AdminContext) {
	logger := reqctx.Logger(c)

	var reqParams filter.SearchTransactions
	if err := c.Bind(&reqParams); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, logger)
		return
	}

	transactions, err := reqctx.Engine(c).GetTransactions(
		c.Request.Context(),
		mappings.MapToMetadata(reqParams.Metadata),
		reqParams.Conditions.ToDbConditions(),
		mappings.MapToQueryParams(reqParams.QueryParams),
	)
	if err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	}

	contracts := make([]*models.Transaction, 0)
	for _, transaction := range transactions {
		contracts = append(contracts, mappings.MapToOldTransactionContractForAdmin(transaction))
	}
}
