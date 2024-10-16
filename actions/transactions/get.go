package transactions

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
)

// get will fetch a transaction
// @Summary		Get transaction by id - Use (GET) /api/v1/transactions/{id} instead.
// @Description	This endpoint has been deprecated. Use (GET) /api/v1/transactions/{id} instead.
// @Tags		Transactions
// @Produce		json
// @Param		id query string true "id"
// @Success		200 {object} models.Transaction "Transaction"
// @Failure		400	"Bad request - Transaction not found or associated with another xpub"
// @Failure 	500	"Internal Server Error - Error while fetching transaction"
// @DeprecatedRouter	/v1/transaction [get]
// @Security	x-auth-xpub
func get(c *gin.Context, userContext *reqctx.UserContext) {
	logger := reqctx.Logger(c)
	id := c.Query("id")

	transaction, err := reqctx.Engine(c).GetTransaction(
		c.Request.Context(),
		userContext.GetXPubID(),
		id,
	)
	if err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	} else if transaction == nil {
		spverrors.ErrorResponse(c, spverrors.ErrCouldNotFindTransaction, logger)
		return
	} else if !transaction.IsXpubIDAssociated(userContext.GetXPubID()) {
		spverrors.ErrorResponse(c, spverrors.ErrAuthorization, logger)
		return
	}

	contract := mappings.MapToOldTransactionContract(transaction)
	c.JSON(http.StatusOK, contract)
}

// getByID will fetch a transaction by id
// Get transaction by id godoc
// @Summary		Get transaction by id
// @Description	Get transaction by id
// @Tags		Transactions
// @Produce		json
// @Param		id path string true "id"
// @Success		200 {object} response.Transaction "Transaction"
// @Failure		400	"Bad request - Transaction not found or associated with another xpub"
// @Failure 	500	"Internal Server Error - Error while fetching transaction"
// @Router		/api/v1/transactions/{id} [get]
// @Security	x-auth-xpub
func getByID(c *gin.Context, userContext *reqctx.UserContext) {
	logger := reqctx.Logger(c)
	id := c.Param("id")

	transaction, err := reqctx.Engine(c).GetTransaction(
		c.Request.Context(),
		userContext.GetXPubID(),
		id,
	)
	if err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	} else if transaction == nil {
		spverrors.ErrorResponse(c, spverrors.ErrCouldNotFindTransaction, logger)
		return
	} else if !transaction.IsXpubIDAssociated(userContext.GetXPubID()) {
		spverrors.ErrorResponse(c, spverrors.ErrAuthorization, logger)
		return
	}

	contract := mappings.MapToTransactionContract(transaction)
	c.JSON(http.StatusOK, contract)
}
