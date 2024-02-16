package transactions

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/actions"
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
)

// search will fetch a list of transactions filtered on conditions and metadata
// Search transaction godoc
// @Summary		Search transaction
// @Description	Search transaction
// @Tags		Transactions
// @Produce		json
// @Param		page query int false "page"
// @Param		page_size query int false "page_size"
// @Param		order_by_field query string false "order_by_field"
// @Param		sort_direction query string false "sort_direction"
// @Param		metadata query string false "metadata"
// @Param		conditions query string false "conditions"
// @Success		200
// @Router		/v1/transaction/search [post]
// @Security	x-auth-xpub
func (a *Action) search(c *gin.Context) {
	reqXPubID := c.GetString(auth.ParamXPubHashKey)

	queryParams, metadata, conditions, err := actions.GetQueryParameters(c)
	if err != nil {
		c.JSON(http.StatusExpectationFailed, err.Error())
		return
	}

	// Record a new transaction (get the hex from parameters)a
	var transactions []*engine.Transaction
	if transactions, err = a.Services.SpvWalletEngine.GetTransactionsByXpubID(
		c.Request.Context(),
		reqXPubID,
		metadata,
		conditions,
		queryParams,
	); err != nil {
		c.JSON(http.StatusExpectationFailed, err.Error())
		return
	}

	contracts := make([]*models.Transaction, 0)
	for _, transaction := range transactions {
		contracts = append(contracts, mappings.MapToTransactionContract(transaction))
	}

	c.JSON(http.StatusOK, contracts)
}
