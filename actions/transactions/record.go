package transactions

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/actions"
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/server/auth"
)

// record will save and complete a transaction
// Record transaction godoc
// @Summary		Record transaction
// @Description	Record transaction
// @Tags		Transactions
// @Produce		json
// @Param		hex query string true "hex"
// @Param		reference_id query string true "reference_id"
// @Param		metadata query string false "metadata"
// @Success		200
// @Router		/v1/transaction/record [post]
// @Security	x-auth-xpub
func (a *Action) record(c *gin.Context) {
	reqXPub := c.GetString(auth.ParamXPubKey)

	var requestBody RecordTransaction
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	xPub, err := a.Services.SpvWalletEngine.GetXpub(c.Request.Context(), reqXPub)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	} else if xPub == nil {
		c.JSON(http.StatusForbidden, actions.ErrXpubNotFound.Error())
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
		requestBody.ReferenceId,
		opts...,
	); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	contract := mappings.MapToTransactionContract(transaction)
	c.JSON(http.StatusCreated, contract)
}
