package transactions

import (
	"github.com/bitcoin-sv/spv-wallet/actions"
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
	"github.com/gin-gonic/gin"
	"net/http"
)

// record will save and complete a transaction
// @Deprecated
// @Summary		Record transaction
// @Description	This endpoint has been deprecated. Use (POST) /api/v1/transactions instead.
// @Tags		Transactions
// @Produce		json
// @Param		RecordTransaction body RecordTransaction true "Transaction to be recorded"
// @Success		201 {object} models.Transaction "Created transaction"
// @Failure		400	"Bad request - Error while parsing RecordTransaction from request body or xpub not found"
// @Failure 	500	"Internal Server Error - Error while recording transaction"
// @Router		/v1/transaction/record [post]
// @Security	x-auth-xpub
func (a *Action) record(c *gin.Context) {
	a.recordTransaction(c)
}

// recordTransaction will save and complete a transaction
// @Summary		Record transaction
// @Description	Record transaction
// @Tags		Transactions
// @Produce		json
// @Param		RecordTransaction body RecordTransaction true "Transaction to be recorded"
// @Success		201 {object} models.Transaction "Created transaction"
// @Failure		400	"Bad request - Error while parsing RecordTransaction from request body or xpub not found"
// @Failure 	500	"Internal Server Error - Error while recording transaction"
// @Router		/api/v1/transactions [post]
// @Security	x-auth-xpub
func (a *Action) recordTransaction(c *gin.Context) {
	reqXPub := c.GetString(auth.ParamXPubKey)

	var requestBody RecordTransaction
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	xPub, err := a.Services.SpvWalletEngine.GetXpub(c.Request.Context(), reqXPub)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	} else if xPub == nil {
		c.JSON(http.StatusBadRequest, actions.ErrXpubNotFound.Error())
		return
	}

	opts := make([]engine.ModelOps, 0)
	if requestBody.Metadata != nil {
		opts = append(opts, engine.WithMetadatas(requestBody.Metadata))
	}

	// Record a new transaction (get the hex from parameters)
	var transaction *engine.Transaction
	if transaction, err = a.Services.SpvWalletEngine.RecordTransaction(
		c.Request.Context(),
		reqXPub,
		requestBody.Hex,
		requestBody.ReferenceID,
		opts...,
	); err != nil {
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	contract := mappings.MapToTransactionContract(transaction)
	c.JSON(http.StatusCreated, contract)
}
