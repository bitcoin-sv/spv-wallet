package transactions

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
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
func record(c *gin.Context, userContext *reqctx.UserContext) {
	logger := reqctx.Logger(c)
	engineInstance := reqctx.Engine(c)

	xpub, err := userContext.ShouldGetXPub()
	if err != nil {
		spverrors.AbortWithErrorResponse(c, err, logger)
		return
	}

	var requestBody OldRecordTransaction
	err = c.ShouldBindJSON(&requestBody)
	if err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	}

	opts := make([]engine.ModelOps, 0)
	if requestBody.Metadata != nil {
		opts = append(opts, engine.WithMetadatas(requestBody.Metadata))
	}

	// Record a new transaction (get the hex from parameters)
	var transaction *engine.Transaction
	if transaction, err = engineInstance.RecordTransaction(
		c.Request.Context(),
		xpub,
		requestBody.Hex,
		requestBody.ReferenceID,
		opts...,
	); err != nil {
		spverrors.ErrorResponse(c, err, logger)
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
func recordTransaction(c *gin.Context, userContext *reqctx.UserContext) {
	logger := reqctx.Logger(c)
	engineInstance := reqctx.Engine(c)

	xpub, err := userContext.ShouldGetXPub()
	if err != nil {
		spverrors.AbortWithErrorResponse(c, err, logger)
		return
	}

	var requestBody RecordTransaction
	err = c.ShouldBindJSON(&requestBody)
	if err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	}

	opts := make([]engine.ModelOps, 0)
	if requestBody.Metadata != nil {
		opts = append(opts, engine.WithMetadatas(requestBody.Metadata))
	}

	// Record a new transaction (get the hex from parameters)
	var transaction *engine.Transaction
	if transaction, err = engineInstance.RecordTransaction(
		c.Request.Context(),
		xpub,
		requestBody.Hex,
		requestBody.ReferenceID,
		opts...,
	); err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	}

	contract := mappings.MapToTransactionContract(transaction)
	c.JSON(http.StatusCreated, contract)
}
