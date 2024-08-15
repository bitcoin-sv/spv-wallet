package transactions

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
	"github.com/gin-gonic/gin"
)

// newTransaction will create a new transaction
// New transaction godoc
// @Summary		New transaction - Use (POST) /api/v1/transactions/drafts instead.
// @Description	This endpoint has been deprecated. Use (POST) /api/v1/transactions/drafts instead.
// @Tags		Transactions
// @Produce		json
// @Param		NewTransaction body NewTransaction true "NewTransaction model containing the transaction config and metadata"
// @Success		201 {object} models.DraftTransaction "Created transaction"
// @Failure		400	"Bad request - Error while parsing NewTransaction from request body or xpub not found"
// @Failure 	500	"Internal Server Error - Error while creating transaction"
// @DeprecatedRouter	/v1/transaction [post]
// @Security	x-auth-xpub
func (a *Action) newTransaction(c *gin.Context) {
	reqXPub := c.GetString(auth.ParamXPubKey)

	xPub, err := a.Services.SpvWalletEngine.GetXpub(c.Request.Context(), reqXPub)
	if err != nil {
		spverrors.ErrorResponse(c, err, a.Services.Logger)
		return
	} else if xPub == nil {
		spverrors.ErrorResponse(c, spverrors.ErrCouldNotFindXpub, a.Services.Logger)
		return
	}

	var requestBody OldNewTransaction
	if err = c.Bind(&requestBody); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, a.Services.Logger)
		return
	}

	opts := a.Services.SpvWalletEngine.DefaultModelOptions()
	if requestBody.Metadata != nil {
		opts = append(opts, engine.WithMetadatas(requestBody.Metadata))
	}

	txConfig := mappings.MapOldTransactionConfigEngineToModel(&requestBody.Config)

	var transaction *engine.DraftTransaction
	if transaction, err = a.Services.SpvWalletEngine.NewTransaction(
		c.Request.Context(),
		xPub.RawXpub(),
		txConfig,
		opts...,
	); err != nil {
		spverrors.ErrorResponse(c, err, a.Services.Logger)
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
// @Param		NewTransaction body NewTransaction true "NewTransaction model containing the transaction config and metadata"
// @Success		201 {object} models.DraftTransaction "Created transaction"
// @Failure		400	"Bad request - Error while parsing NewTransaction from request body or xpub not found"
// @Failure 	500	"Internal Server Error - Error while creating transaction"
// @Router		/api/v1/transactions/drafts [post]
// @Security	x-auth-xpub
func (a *Action) newTransactionDraft(c *gin.Context) {
	reqXPub := c.GetString(auth.ParamXPubKey)

	xPub, err := a.Services.SpvWalletEngine.GetXpub(c.Request.Context(), reqXPub)
	if err != nil {
		spverrors.ErrorResponse(c, err, a.Services.Logger)
		return
	} else if xPub == nil {
		spverrors.ErrorResponse(c, spverrors.ErrCouldNotFindXpub, a.Services.Logger)
		return
	}

	var requestBody NewTransaction
	if err = c.Bind(&requestBody); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, a.Services.Logger)
		return
	}

	opts := a.Services.SpvWalletEngine.DefaultModelOptions()
	if requestBody.Metadata != nil {
		opts = append(opts, engine.WithMetadatas(requestBody.Metadata))
	}

	txConfig := mappings.MapTransactionConfigEngineToModel(&requestBody.Config)

	var transaction *engine.DraftTransaction
	if transaction, err = a.Services.SpvWalletEngine.NewTransaction(
		c.Request.Context(),
		xPub.RawXpub(),
		txConfig,
		opts...,
	); err != nil {
		spverrors.ErrorResponse(c, err, a.Services.Logger)
		return
	}

	contract := mappings.MapToDraftTransactionContract(transaction)
	c.JSON(http.StatusCreated, contract)
}
