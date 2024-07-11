package transactions

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
	"github.com/gin-gonic/gin"
)

// update will update a transaction
// Update transaction godoc
// @Summary		Update transaction - Use (PATCH) /api/v1/transactions/{id} instead.
// @Description	This endpoint has been deprecated. Use (PATCH) /api/v1/transactions/{id} instead.
// @Tags		Transactions
// @Produce		json
// @Param		UpdateTransaction body UpdateTransaction true " "
// @Success		200 {object} models.Transaction "Updated transaction"
// @Failure		400	"Bad request - Error while parsing UpdateTransaction from request body, tx not found or tx is not associated with the xpub"
// @Failure 	500	"Internal Server Error - Error while updating transaction"
// @Router		/v1/transaction [patch]
// @Security	x-auth-xpub
// @Deprecated
func (a *Action) update(c *gin.Context) {

	var requestBody UpdateTransaction
	if err := c.Bind(&requestBody); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, a.Services.Logger)
		return
	}
	id := requestBody.ID

	a.updateTransactionWithID(c, id, requestBody.Metadata)
}

// update will update a transaction
// Update transaction godoc
// @Summary		Update transaction
// @Description	Update transaction
// @Tags		Transactions
// @Produce		json
// @Param		UpdateTransactionRequest body UpdateTransactionRequest true " "
// @Success		200 {object} models.Transaction "Updated transaction"
// @Failure		400	"Bad request - Error while parsing UpdateTransaction from request body, tx not found or tx is not associated with the xpub"
// @Failure 	500	"Internal Server Error - Error while updating transaction"
// @Router		/api/v1/transactions/{id} [patch]
// @Security	x-auth-xpub
func (a *Action) updateTransaction(c *gin.Context) {
	var requestBody UpdateTransactionRequest
	if err := c.Bind(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	id := c.Param("id")

	a.updateTransactionWithID(c, id, requestBody.Metadata)
}

func (a *Action) updateTransactionWithID(c *gin.Context, id string, requestMetadata engine.Metadata) {
	reqXPubID := c.GetString(auth.ParamXPubHashKey)

	// Get a transaction by ID
	transaction, err := a.Services.SpvWalletEngine.UpdateTransactionMetadata(
		c.Request.Context(),
		reqXPubID,
		id,
		requestMetadata,
	)
	if err != nil {
		spverrors.ErrorResponse(c, err, a.Services.Logger)
		return
	} else if transaction == nil {
		spverrors.ErrorResponse(c, spverrors.ErrCouldNotFindTransaction, a.Services.Logger)
	} else if !transaction.IsXpubIDAssociated(reqXPubID) {
		spverrors.ErrorResponse(c, spverrors.ErrAuthorization, a.Services.Logger)
		return
	}

	contract := mappings.MapToTransactionContract(transaction)
	c.JSON(http.StatusOK, contract)
}
