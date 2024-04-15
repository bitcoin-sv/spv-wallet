package transactions

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
	"github.com/gin-gonic/gin"
)

// search will fetch a list of transactions filtered on conditions and metadata
// Search transaction godoc
// @Summary		Search transaction
// @Description	Search transaction
// @Tags		Transactions
// @Produce		json
// @Param		SearchTransactions body SearchTransactions false "Supports targeted resource searches with filters for metadata and custom conditions, plus options for pagination and sorting to streamline data exploration and analysis"
// @Success		200 {object} []models.Transaction "List of transactions"
// @Failure		400	"Bad request - Error while parsing SearchTransactions from request body"
// @Failure 	500	"Internal server error - Error while searching for transactions"
// @Router		/v1/transaction/search [post]
// @Security	x-auth-xpub
func (a *Action) search(c *gin.Context) {
	reqXPubID := c.GetString(auth.ParamXPubHashKey)

	var reqParams SearchTransactions
	if err := c.Bind(&reqParams); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// Record a new transaction (get the hex from parameters)a
	transactions, err := a.Services.SpvWalletEngine.GetTransactionsByXpubID(
		c.Request.Context(),
		reqXPubID,
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
		contracts = append(contracts, mappings.MapToTransactionContract(transaction))
	}

	c.JSON(http.StatusOK, contracts)
}
