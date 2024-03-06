package admin

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/actions"
	"github.com/bitcoin-sv/spv-wallet/engine"
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
// @Param		SearchRequestParameters body actions.SearchRequestParameters false "SearchRequestParameters model containing metadata, conditions and query params"
// @Success		200
// @Router		/v1/admin/transactions/search [post]
// @Security	x-auth-xpub
func (a *Action) transactionsSearch(c *gin.Context) {
	queryParams, metadata, conditions, err := actions.GetSearchQueryParameters(c)
	if err != nil {
		c.JSON(http.StatusExpectationFailed, err.Error())
		return
	}

	var transactions []*engine.Transaction
	if transactions, err = a.Services.SpvWalletEngine.GetTransactions(
		c.Request.Context(),
		metadata,
		conditions,
		queryParams,
	); err != nil {
		c.JSON(http.StatusExpectationFailed, err.Error())
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
// @Param		CountRequestParameters body actions.CountRequestParameters false "CountRequestParameters model containing metadata and conditions"
// @Success		200
// @Router		/v1/admin/transactions/count [post]
// @Security	x-auth-xpub
func (a *Action) transactionsCount(c *gin.Context) {
	metadata, conditions, err := actions.GetCountQueryParameters(c)
	if err != nil {
		c.JSON(http.StatusExpectationFailed, err.Error())
		return
	}

	var count int64
	if count, err = a.Services.SpvWalletEngine.GetTransactionsCount(
		c.Request.Context(),
		metadata,
		conditions,
	); err != nil {
		c.JSON(http.StatusExpectationFailed, err.Error())
		return
	}

	c.JSON(http.StatusOK, count)
}
