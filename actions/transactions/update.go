package transactions

import (
	spverrors2 "github.com/bitcoin-sv/spv-wallet/engine/spverrors"
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
	reqXPubID := c.GetString(auth.ParamXPubHashKey)

	var requestBody UpdateTransaction
	if err := c.Bind(&requestBody); err != nil {
		spverrors2.ErrorResponse(c, spverrors2.ErrCannotBindRequest, a.Services.Logger)
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
		spverrors2.ErrorResponse(c, err, a.Services.Logger)
		return
	} else if transaction == nil {
		spverrors2.ErrorResponse(c, spverrors2.ErrCouldNotFindTransaction, a.Services.Logger)
	} else if !transaction.IsXpubIDAssociated(reqXPubID) {
		spverrors2.ErrorResponse(c, spverrors2.ErrAuthorization, a.Services.Logger)
		return
	}

	contract := mappings.MapToTransactionContract(transaction)
	c.JSON(http.StatusOK, contract)
}
