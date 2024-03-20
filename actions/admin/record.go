package admin

import (
	"errors"
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/gin-gonic/gin"
	"github.com/mrz1836/go-datastore"
)

// transactionRecord will save and complete a transaction directly, without any checks
// Record transactions godoc
// @Summary		Record transactions
// @Description	Record transactions
// @Tags		Admin
// @Produce		json
// @Param		RecordTransaction body RecordTransaction true "RecordTransaction model containing hex of the transaction to record"
// @Success		201	{object} models.Transaction "Recorded transaction"
// @Failure		400	"Bad request - Error while parsing RecordTransaction from request body"
// @Failure 	500	"Internal Server Error - Error while fetching count of access keys"
// @Router		/v1/admin/transactions/record [post]
// @Security	x-auth-xpub
func (a *Action) transactionRecord(c *gin.Context) {
	var requestBody RecordTransaction
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
				c.JSON(http.StatusInternalServerError, err.Error())
				return
			}
		} else {
			c.JSON(http.StatusInternalServerError, err.Error())
			return
		}
	}

	contract := mappings.MapToTransactionContract(transaction)

	c.JSON(http.StatusCreated, contract)
}
