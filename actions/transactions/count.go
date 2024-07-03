package transactions

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
	"github.com/gin-gonic/gin"
)

// count will fetch a count of transactions filtered on conditions and metadata
// Count of transactions godoc
// @Summary		Count of transactions
// @Description	The functionality of this method will be offered in the future by /api/v1/transactions [get].
// @Tags		Transactions
// @Produce		json
// @Param		CountTransactions body filter.CountTransactions false "Enables filtering of elements to be counted"
// @Success		200	{number} int64 "Count of access keys"
// @Failure		400	"Bad request - Error while parsing CountTransactions from request body"
// @Failure 	500	"Internal Server Error - Error while fetching count of transactions"
// @Router		/v1/transaction/count [post]
// @Security	x-auth-xpub
func (a *Action) count(c *gin.Context) {
	reqXPubID := c.GetString(auth.ParamXPubHashKey)

	var reqParams filter.CountTransactions
	if err := c.Bind(&reqParams); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	count, err := a.Services.SpvWalletEngine.GetTransactionsByXpubIDCount(
		c.Request.Context(),
		reqXPubID,
		mappings.MapToMetadata(reqParams.Metadata),
		reqParams.Conditions.ToDbConditions(),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, count)
}
