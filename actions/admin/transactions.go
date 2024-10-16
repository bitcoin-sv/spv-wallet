package admin

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
)

// transactionsSearch will fetch a list of transactions filtered by metadata
// Search for transactions filtering by metadata godoc
// @Summary		Search for transactions
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

// getAdminByID fetches a transaction by id for admins
// @Summary		Get transaction by id for admins
// @Description	Get transaction by id for admins
// @Tags		Admin
// @Produce		json
// @Param		id path string true "Transaction ID"
// @Success		200 {object} models.Transaction "Transaction"
// @Failure		400	"Bad request - Transaction not found or error in data fetching"
// @Failure 	500	"Internal Server Error - Error while fetching transaction"
// @Router		/v1/admin/transactions/{id} [get]
// @Security	x-auth-xpub
func getAdminByID(c *gin.Context, adminContext *reqctx.AdminContext) {
	logger := reqctx.Logger(c)
	id := c.Param("id")

	transaction, err := reqctx.Engine(c).GetAdminTransaction(c, id)
	if err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	}
	if transaction == nil {
		spverrors.ErrorResponse(c, spverrors.ErrCouldNotFindTransaction, logger)
		return
	}

	contract := mappings.MapToTransactionContract(transaction)
	c.JSON(http.StatusOK, contract)
}

// getAdminDeprecated fetches a transaction for admins using a deprecated path
// @Deprecated This endpoint has been deprecated. Use /api/v1/admin/transactions/{id} instead.
// @Summary		Get transaction by id for admins (Deprecated)
// @Description	This endpoint has been deprecated. Use /api/v1/admin/transactions/{id} instead.
// @Tags		Admin
// @Produce		json
// @Param		id query string true "id"
// @Deprecated
// @Success		200 {object} models.Transaction "Transaction"
// @Failure		400	"Bad request - Transaction not found or error in data fetching"
// @Failure 	500	"Internal Server Error - Error while fetching transaction"
// @Router		/v1/admin/transactions [get]
// @Security	x-auth-xpub
func getAdminDeprecated(c *gin.Context, adminContext *reqctx.AdminContext) {
	logger := reqctx.Logger(c)
	id := c.Query("id")

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
