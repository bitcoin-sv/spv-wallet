package transactions

import (
	spverrors2 "github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
	"github.com/gin-gonic/gin"
)

// newTransaction will create a new transaction
// New transaction godoc
// @Summary		New transaction
// @Description	New transaction
// @Tags		Transactions
// @Produce		json
// @Param		NewTransaction body NewTransaction true "NewTransaction model containing the transaction config and metadata"
// @Success		201 {object} models.DraftTransaction "Created transaction"
// @Failure		400	"Bad request - Error while parsing NewTransaction from request body or xpub not found"
// @Failure 	500	"Internal Server Error - Error while creating transaction"
// @Router		/v1/transaction [post]
// @Security	x-auth-xpub
func (a *Action) newTransaction(c *gin.Context) {
	reqXPub := c.GetString(auth.ParamXPubKey)

	xPub, err := a.Services.SpvWalletEngine.GetXpub(c.Request.Context(), reqXPub)
	if err != nil {
		spverrors2.ErrorResponse(c, err, a.Services.Logger)
		return
	} else if xPub == nil {
		spverrors2.ErrorResponse(c, spverrors2.ErrCouldNotFindXpub, a.Services.Logger)
		return
	}

	var requestBody NewTransaction
	if err = c.Bind(&requestBody); err != nil {
		spverrors2.ErrorResponse(c, spverrors2.ErrCannotBindRequest, a.Services.Logger)
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
		spverrors2.ErrorResponse(c, err, a.Services.Logger)
		return
	}

	contract := mappings.MapToDraftTransactionContract(transaction)
	c.JSON(http.StatusCreated, contract)
}
