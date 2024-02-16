package transactions

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
)

// get will fetch a transaction
// Get transaction by id godoc
// @Summary		Get transaction by id
// @Description	Get transaction by id
// @Tags		Transactions
// @Produce		json
// @Param		id query string true "id"
// @Success		200
// @Router		/v1/transaction [get]
// @Security	x-auth-xpub
func (a *Action) get(c *gin.Context) {
	reqXPubID := c.GetString(auth.ParamXPubHashKey)
	id := c.Query("id")

	transaction, err := a.Services.SpvWalletEngine.GetTransaction(
		c.Request.Context(),
		reqXPubID,
		id,
	)
	if err != nil {
		c.JSON(http.StatusExpectationFailed, err.Error())
		return
	} else if transaction == nil {
		c.JSON(http.StatusNotFound, "not found")
		return
	} else if !transaction.IsXpubIDAssociated(reqXPubID) {
		c.JSON(http.StatusForbidden, "unauthorized")
		return
	}

	contract := mappings.MapToTransactionContract(transaction)
	c.JSON(http.StatusOK, contract)
}
