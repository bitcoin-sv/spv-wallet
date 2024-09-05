package transactions

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

// newTransaction will create a new transaction
// New transaction godoc
// @Summary		New transaction - Use (POST) /api/v1/transactions/drafts instead.
// @Description	This endpoint has been deprecated. Use (POST) /api/v1/transactions/drafts instead.
// @Tags		Transactions
// @Produce		json
// @Param		OldNewDraftTransaction body OldNewDraftTransaction true "OldNewDraftTransaction model containing the transaction config and metadata"
// @Success		201 {object} OldNewDraftTransaction "Created transaction"
// @Failure		400	"Bad request - Error while parsing OldNewDraftTransaction from request body or xpub not found"
// @Failure 	500	"Internal Server Error - Error while creating transaction"
// @DeprecatedRouter	/v1/transaction [post]
// @Security	x-auth-xpub
func newTransaction(c *gin.Context, userContext *reqctx.UserContext) {
	logger := reqctx.Logger(c)
	engineInstance := reqctx.Engine(c)

	xpub, err := userContext.ShouldGetXPub()
	if err != nil {
		spverrors.AbortWithErrorResponse(c, err, logger)
		return
	}

	var requestBody OldNewDraftTransaction
	err = c.Bind(&requestBody)
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, logger)
		return
	}

	opts := engineInstance.DefaultModelOptions()
	if requestBody.Metadata != nil {
		opts = append(opts, engine.WithMetadatas(requestBody.Metadata))
	}

	txConfig := mappings.MapOldTransactionConfigEngineToModel(&requestBody.Config)

	var transaction *engine.DraftTransaction
	if transaction, err = engineInstance.NewTransaction(
		c.Request.Context(),
		xpub,
		txConfig,
		opts...,
	); err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	}

	contract := mappings.MapToOldDraftTransactionContract(transaction)
	c.JSON(http.StatusCreated, contract)
}

// newTransactionDraft will create a new transaction draft
// New transaction draft godoc
// @Summary		New transaction draft
// @Description	New transaction draft
// @Tags		Transactions
// @Produce		json
// @Param		NewDraftTransaction body NewDraftTransaction true "NewDraftTransaction model containing the transaction config and metadata"
// @Success		201 {object} response.DraftTransaction "Created transaction"
// @Failure		400	"Bad request - Error while parsing NewDraftTransaction from request body or xpub not found"
// @Failure 	500	"Internal Server Error - Error while creating transaction"
// @Router		/api/v1/transactions/drafts [post]
// @Security	x-auth-xpub
func newTransactionDraft(c *gin.Context, userContext *reqctx.UserContext) {
	logger := reqctx.Logger(c)
	engineInstance := reqctx.Engine(c)

	xpub, err := userContext.ShouldGetXPub()
	if err != nil {
		spverrors.AbortWithErrorResponse(c, err, logger)
		return
	}

	var requestBody NewDraftTransaction
	err = c.Bind(&requestBody)
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, logger)
		return
	}

	opts := engineInstance.DefaultModelOptions()
	if requestBody.Metadata != nil {
		opts = append(opts, engine.WithMetadatas(requestBody.Metadata))
	}

	txConfig := mappings.MapTransactionConfigEngineToModel(&requestBody.Config)

	var transaction *engine.DraftTransaction
	if transaction, err = engineInstance.NewTransaction(
		c.Request.Context(),
		xpub,
		txConfig,
		opts...,
	); err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	}

	contract := mappings.MapToDraftTransactionContract(transaction)
	c.JSON(http.StatusCreated, contract)
}
