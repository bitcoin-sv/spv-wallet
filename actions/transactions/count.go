package transactions

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/actions"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
	"github.com/gin-gonic/gin"
)

// count will fetch a count of transactions filtered on conditions and metadata
// Count of transactions godoc
// @Summary		Count of transactions
// @Description	Count of transactions
// @Tags		Transactions
// @Produce		json
// @Param		CountRequestParameters body actions.CountRequestParameters false "CountRequestParameters model containing metadata and conditions"
// @Success		200
// @Router		/v1/transaction/count [post]
// @Security	x-auth-xpub
func (a *Action) count(c *gin.Context) {
	reqXPubID := c.GetString(auth.ParamXPubHashKey)

	metadata, conditions, err := actions.GetCountQueryParameters(c)
	if err != nil {
		c.JSON(http.StatusExpectationFailed, err.Error())
		return
	}

	var count int64
	if count, err = a.Services.SpvWalletEngine.GetTransactionsByXpubIDCount(
		c.Request.Context(),
		reqXPubID,
		metadata,
		conditions,
	); err != nil {
		c.JSON(http.StatusExpectationFailed, err.Error())
		return
	}

	c.JSON(http.StatusOK, count)
}
