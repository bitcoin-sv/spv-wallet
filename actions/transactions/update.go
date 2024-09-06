package transactions

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

// update will update a transaction
// Update transaction godoc
// @Summary		Update transaction - Use (PATCH) /api/v1/transactions/{id} instead.
// @Description	This endpoint has been deprecated. Use (PATCH) /api/v1/transactions/{id} instead.
// @Tags		Transactions
// @Produce		json
// @Param		UpdateTransaction body UpdateTransaction true "Pass update transaction request model in the body"
// @Success		200 {object} models.Transaction "Updated transaction"
// @Failure		400	"Bad request - Error while parsing UpdateTransaction from request body, tx not found or tx is not associated with the xpub"
// @Failure 	500	"Internal Server Error - Error while updating transaction"
// @DeprecatedRouter	/v1/transaction [patch]
// @Security	x-auth-xpub
func update(c *gin.Context, userContext *reqctx.UserContext) {
	logger := reqctx.Logger(c)

	var requestBody OldUpdateTransaction
	if err := c.Bind(&requestBody); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, logger)
		return
	}
	id := requestBody.ID

	// Get a transaction by ID
	transaction, err := reqctx.Engine(c).UpdateTransactionMetadata(
		c.Request.Context(),
		userContext.GetXPubID(),
		id,
		requestBody.Metadata,
	)
	if err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	} else if transaction == nil {
		spverrors.ErrorResponse(c, spverrors.ErrCouldNotFindTransaction, logger)
	} else if !transaction.IsXpubIDAssociated(userContext.GetXPubID()) {
		spverrors.ErrorResponse(c, spverrors.ErrAuthorization, logger)
		return
	}

	contract := mappings.MapToOldTransactionContract(transaction)
	c.JSON(http.StatusOK, contract)
}

// update will update a transaction metadata
// Update transaction godoc
// @Summary		Update transaction metadata
// @Description	Update transaction metadata
// @Tags		Transactions
// @Produce		json
// @Param		UpdateTransactionRequest body UpdateTransactionRequest true "Pass update transaction request model in the body with updated metadata"
// @Success		200 {object} response.Transaction "Updated transaction metadata"
// @Failure		400	"Bad request - Error while parsing UpdateTransaction from request body, tx not found or tx is not associated with the xpub"
// @Failure 	500	"Internal Server Error - Error while updating transaction metadata"
// @Router		/api/v1/transactions/{id} [patch]
// @Security	x-auth-xpub
func updateTransactionMetadata(c *gin.Context, userContext *reqctx.UserContext) {
	logger := reqctx.Logger(c)
	var requestBody UpdateTransactionRequest
	if err := c.Bind(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	id := c.Param("id")

	// Get a transaction by ID
	transaction, err := reqctx.Engine(c).UpdateTransactionMetadata(
		c.Request.Context(),
		userContext.GetXPubID(),
		id,
		requestBody.Metadata,
	)
	if err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	} else if transaction == nil {
		spverrors.ErrorResponse(c, spverrors.ErrCouldNotFindTransaction, logger)
	} else if !transaction.IsXpubIDAssociated(userContext.GetXPubID()) {
		spverrors.ErrorResponse(c, spverrors.ErrAuthorization, logger)
		return
	}

	contract := mappings.MapToTransactionContract(transaction)
	c.JSON(http.StatusOK, contract)
}
