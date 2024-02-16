package admin

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/mrz1836/go-datastore"
)

// transactionRecord will save and complete a transaction directly, without any checks
// Record transactions godoc
// @Summary		Record transactions
// @Description	Record transactions
// @Tags		Admin
// @Produce		json
// @Param		hex query string true "Transaction hex"
// @Success		201
// @Router		/v1/admin/transactions/record [post]
// @Security	x-auth-xpub
func (a *Action) transactionRecord(c *gin.Context) {
	var requestBody AdminRecordTransaction
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	// Set the metadata
	opts := make([]engine.ModelOps, 0)

	// Record a new transaction (get the hex from parameters)
	transaction, err := a.Services.SpvWalletEngine.RecordRawTransaction(
		c.Request.Context(),
		requestBody.Hex,
		opts...,
	)
	if err != nil {
		if errors.Is(err, datastore.ErrDuplicateKey) {
			// already registered, just return the registered transaction
			if transaction, err = a.Services.SpvWalletEngine.GetTransactionByHex(c.Request.Context(), requestBody.Hex); err != nil {
				c.JSON(http.StatusUnprocessableEntity, err.Error())
				return
			}
		} else {
			c.JSON(http.StatusUnprocessableEntity, err.Error())
			return
		}
	}

	contract := mappings.MapToTransactionContract(transaction)

	c.JSON(http.StatusCreated, contract)
}
