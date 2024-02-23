package transactions

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
	apirouter "github.com/mrz1836/go-api-router"
)

// update will update a transaction
// Update transaction godoc
// @Summary		Update transaction
// @Description	Update transaction
// @Tags		Transactions
// @Produce		json
// @Param		id query string true "id"
// @Param		metadata query string true "metadata"
// @Success		200
// @Router		/v1/transaction [patch]
// @Security	x-auth-xpub
func (a *Action) update(c *gin.Context) {
	reqXPubID := c.GetString(auth.ParamXPubHashKey)

	var requestBody UpdateTransaction
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		apirouter.ReturnResponse(c.Writer, c.Request, http.StatusExpectationFailed, err.Error())
		return
	}

	// Get a transaction by ID
	transaction, err := a.Services.SpvWalletEngine.UpdateTransactionMetadata(
		c.Request.Context(),
		reqXPubID,
		requestBody.ID,
		requestBody.Metadata,
	)
	if err != nil {
		c.JSON(http.StatusExpectationFailed, err.Error())
		return
	} else if transaction == nil {
		c.JSON(http.StatusNotFound, "not found")
	} else if !transaction.IsXpubIDAssociated(reqXPubID) {
		c.JSON(http.StatusForbidden, "unauthorized")
		return
	}

	contract := mappings.MapToTransactionContract(transaction)
	c.JSON(http.StatusOK, contract)
}
