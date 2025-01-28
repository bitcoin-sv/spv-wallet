package transactions

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

// update will update a transaction metadata
// Update transaction godoc
// @Summary		Update transaction metadata
// @Description	Update transaction metadata
// @Tags		Transactions
// @Produce		json
// @Param		id path string true "id"
// @Param		UpdateTransactionRequest body UpdateTransactionRequest true "Pass update transaction request model in the body with updated metadata"
// @Success		200 {object} response.Transaction "Updated transaction metadata"
// @Failure		400	"Bad request - Error while parsing UpdateTransactionRequest from request body, tx not found or tx is not associated with the xpub"
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
