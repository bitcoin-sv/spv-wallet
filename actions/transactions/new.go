package transactions

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/actions"
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
)

// newTransaction will create a new transaction
// New transaction godoc
// @Summary		New transaction
// @Description	New transaction
// @Tags		Transactions
// @Produce		json
// @Param		config query string true "transaction config"
// @Param		metadata query string false "metadata"
// @Success		201
// @Router		/v1/transaction [post]
// @Security	x-auth-xpub
func (a *Action) newTransaction(c *gin.Context) {
	reqXPub := c.GetString(auth.ParamXPubKey)

	xPub, err := a.Services.SpvWalletEngine.GetXpub(c.Request.Context(), reqXPub)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	} else if xPub == nil {
		c.JSON(http.StatusForbidden, actions.ErrXpubNotFound.Error())
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
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	contract := mappings.MapToDraftTransactionContract(transaction)
	c.JSON(http.StatusCreated, contract)
}
