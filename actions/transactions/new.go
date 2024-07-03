package transactions

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/actions"
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
	"github.com/gin-gonic/gin"
)

// newTransaction will create a new transaction
// New transaction godoc
// @Summary		New transaction
// @Description	This endpoint has been deprecated. Use (POST) /api/v1/transactions/drafts instead.
// @Tags		Transactions
// @Produce		json
// @Param		NewTransaction body NewTransaction true "NewTransaction model containing the transaction config and metadata"
// @Success		201 {object} models.DraftTransaction "Created transaction"
// @Failure		400	"Bad request - Error while parsing NewTransaction from request body or xpub not found"
// @Failure 	500	"Internal Server Error - Error while creating transaction"
// @Router		/v1/transaction [post]
// @Security	x-auth-xpub
// @Deprecated
func (a *Action) newTransaction(c *gin.Context) {
	a.newTransactionDraft(c)
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
		c.JSON(http.StatusBadRequest, err.Error())
		return
	} else if xPub == nil {
		c.JSON(http.StatusBadRequest, actions.ErrXpubNotFound.Error())
		return
	}

	var requestBody NewTransaction
	if err = c.Bind(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
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
		c.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	contract := mappings.MapToDraftTransactionContract(transaction)
	c.JSON(http.StatusCreated, contract)
}
