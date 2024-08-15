package transactions

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
	"github.com/gin-gonic/gin"
)

// record will save and complete a transaction
// @Summary		Record transaction - Use (POST) /api/v1/transactions instead.
// @Description	This endpoint has been deprecated. Use (POST) /api/v1/transactions instead.
// @Tags		Transactions
// @Produce		json
// @Param		RecordTransaction body RecordTransaction true "Transaction to be recorded"
// @Success		201 {object} models.Transaction "Created transaction"
// @Failure		400	"Bad request - Error while parsing RecordTransaction from request body or xpub not found"
// @Failure 	500	"Internal Server Error - Error while recording transaction"
// @DeprecatedRouter	/v1/transaction/record [post]
// @Security	x-auth-xpub
func (a *Action) record(c *gin.Context) {
	reqXPub := c.GetString(auth.ParamXPubKey)

	var requestBody OldRecordTransaction
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		spverrors.ErrorResponse(c, err, a.Services.Logger)
		return
	}

	xPub, err := a.Services.SpvWalletEngine.GetXpub(c.Request.Context(), reqXPub)
	if err != nil {
		spverrors.ErrorResponse(c, err, a.Services.Logger)
		return
	} else if xPub == nil {
		spverrors.ErrorResponse(c, spverrors.ErrCouldNotFindXpub, a.Services.Logger)
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
		spverrors.ErrorResponse(c, err, a.Services.Logger)
		return
	}

	contract := mappings.MapToOldTransactionContract(transaction)
	c.JSON(http.StatusCreated, contract)
}

// recordTransaction will save and complete a transaction
// @Summary		Record transaction
// @Description	Record transaction
// @Tags		Transactions
// @Produce		json
// @Param		RecordTransaction body RecordTransaction true "Transaction to be recorded"
// @Success		201 {object} response.Transaction "Created transaction"
// @Failure		400	"Bad request - Error while parsing RecordTransaction from request body or xpub not found"
// @Failure 	500	"Internal Server Error - Error while recording transaction"
// @Router		/api/v1/transactions [post]
// @Security	x-auth-xpub
func (a *Action) recordTransaction(c *gin.Context) {
	reqXPub := c.GetString(auth.ParamXPubKey)

	var requestBody OldRecordTransaction
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		spverrors.ErrorResponse(c, err, a.Services.Logger)
		return
	}

	xPub, err := a.Services.SpvWalletEngine.GetXpub(c.Request.Context(), reqXPub)
	if err != nil {
		spverrors.ErrorResponse(c, err, a.Services.Logger)
		return
	} else if xPub == nil {
		spverrors.ErrorResponse(c, spverrors.ErrCouldNotFindXpub, a.Services.Logger)
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
		spverrors.ErrorResponse(c, err, a.Services.Logger)
		return
	}

	contract := mappings.MapToTransactionContract(transaction)
	c.JSON(http.StatusCreated, contract)
}
