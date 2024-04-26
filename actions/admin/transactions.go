package admin

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/gin-gonic/gin"
)

// transactionsSearch will fetch a list of transactions filtered by metadata
// Search for transactions filtering by metadata godoc
// @Summary		Search for transactions
// @Description	Search for transactions
// @Tags		Admin
// @Produce		json
// @Param		SearchRequestParameters body actions.SearchRequestParameters false "Supports targeted resource searches with filters for metadata and custom conditions, plus options for pagination and sorting to streamline data exploration and analysis"
// @Success		200 {object} []models.Transaction "List of transactions"
// @Failure		400	"Bad request - Error while parsing SearchRequestParameters from request body"
// @Failure 	500	"Internal server error - Error while searching for transactions"
// @Router		/v1/admin/transactions/search [post]
// @Security	x-auth-xpub
func (a *Action) transactionsSearch(c *gin.Context) {
	var reqParams SearchTransactions
	if err := c.Bind(&reqParams); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	transactions, err := a.Services.SpvWalletEngine.GetTransactions(
		c.Request.Context(),
		reqParams.Metadata,
		reqParams.Conditions.ToDbConditions(),
		reqParams.QueryParams,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	contracts := make([]*models.Transaction, 0)
	for _, transaction := range transactions {
		contracts = append(contracts, mappings.MapToTransactionContractForAdmin(transaction))
	}

	c.JSON(http.StatusOK, contracts)
}

// transactionsCount will count all transactions filtered by metadata
// Count transactions filtering by metadata godoc
// @Summary		Count transactions
// @Description	Count transactions
// @Tags		Admin
// @Produce		json
// @Param		CountRequestParameters body actions.CountRequestParameters false "Enables precise filtering of resource counts using custom conditions or metadata, catering to specific business or analysis needs"
// @Success		200	{number} int64 "Count of transactions"
// @Failure		400	"Bad request - Error while parsing CountRequestParameters from request body"
// @Failure 	500	"Internal Server Error - Error while fetching count of transactions"
// @Router		/v1/admin/transactions/count [post]
// @Security	x-auth-xpub
func (a *Action) transactionsCount(c *gin.Context) {
	var reqParams CountTransactions
	if err := c.Bind(&reqParams); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	count, err := a.Services.SpvWalletEngine.GetTransactionsCount(
		c.Request.Context(),
		reqParams.Metadata,
		reqParams.Conditions.ToDbConditions(),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, count)
}
