package transactions

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
	"github.com/gin-gonic/gin"
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
func (a *Action) get(c *gin.Context) {
	id := c.Query("id")
	reqXPubID := c.GetString(auth.ParamXPubHashKey)

	transaction, err := a.Services.SpvWalletEngine.GetTransaction(
		c.Request.Context(),
		reqXPubID,
		id,
	)
	if err != nil {
		spverrors.ErrorResponse(c, err, a.Services.Logger)
		return
	} else if transaction == nil {
		spverrors.ErrorResponse(c, spverrors.ErrCouldNotFindTransaction, a.Services.Logger)
		return
	} else if !transaction.IsXpubIDAssociated(reqXPubID) {
		spverrors.ErrorResponse(c, spverrors.ErrAuthorization, a.Services.Logger)
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
// @Success		200 {object} models.Transaction "Transaction"
// @Failure		400	"Bad request - Transaction not found or associated with another xpub"
// @Failure 	500	"Internal Server Error - Error while fetching transaction"
// @Router		/api/v1/transactions/{id} [get]
// @Security	x-auth-xpub
func (a *Action) getByID(c *gin.Context) {
	id := c.Param("id")
	reqXPubID := c.GetString(auth.ParamXPubHashKey)

	transaction, err := a.Services.SpvWalletEngine.GetTransaction(
		c.Request.Context(),
		reqXPubID,
		id,
	)
	if err != nil {
		spverrors.ErrorResponse(c, err, a.Services.Logger)
		return
	} else if transaction == nil {
		spverrors.ErrorResponse(c, spverrors.ErrCouldNotFindTransaction, a.Services.Logger)
		return
	} else if !transaction.IsXpubIDAssociated(reqXPubID) {
		spverrors.ErrorResponse(c, spverrors.ErrAuthorization, a.Services.Logger)
		return
	}

	contract := mappings.MapToTransactionContract(transaction)
	c.JSON(http.StatusOK, contract)
}
