package transactions

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
	"github.com/gin-gonic/gin"
)

// update will update a transaction
// Update transaction godoc
// @Summary		Update transaction
// @Description	Update transaction
// @Tags		Transactions
// @Produce		json
// @Param		UpdateTransaction body UpdateTransaction true " "
// @Success		200 {object} models.Transaction "Updated transaction"
// @Failure		400	"Bad request - Error while parsing UpdateTransaction from request body, tx not found or tx is not associated with the xpub"
// @Failure 	500	"Internal Server Error - Error while updating transaction"
// @Router		/v1/transaction [patch]
// @Security	x-auth-xpub
func (a *Action) update(c *gin.Context) {

	var requestBody UpdateTransaction
	if err := c.Bind(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	id := requestBody.ID

	a.updateTransactionWithId(c, id, requestBody)
}

func (a *Action) updateTransactionWithId(c *gin.Context, id string, requestBody UpdateTransaction) {
	reqXPubID := c.GetString(auth.ParamXPubHashKey)

	// Get a transaction by ID
	transaction, err := a.Services.SpvWalletEngine.UpdateTransactionMetadata(
		c.Request.Context(),
		reqXPubID,
		id,
		requestBody.Metadata,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	} else if transaction == nil {
		c.JSON(http.StatusBadRequest, "not found")
	} else if !transaction.IsXpubIDAssociated(reqXPubID) {
		c.JSON(http.StatusBadRequest, "unauthorized")
		return
	}

	contract := mappings.MapToTransactionContract(transaction)
	c.JSON(http.StatusOK, contract)
}

// update will update a transaction
// Update transaction godoc
// @Summary		Update transaction
// @Description	Update transaction
// @Tags		New Transactions
// @Produce		json
// @Param		UpdateTransaction body UpdateTransaction true " "
// @Success		200 {object} models.Transaction "Updated transaction"
// @Failure		400	"Bad request - Error while parsing UpdateTransaction from request body, tx not found or tx is not associated with the xpub"
// @Failure 	500	"Internal Server Error - Error while updating transaction"
// @Router		/v1/transactions/{id} [patch]
// @Security	x-auth-xpub
func (a *Action) updateTransaction(c *gin.Context) {
	var requestBody UpdateTransaction
	if err := c.Bind(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	id := c.Param("id")

	a.updateTransactionWithId(c, id, requestBody)
}
