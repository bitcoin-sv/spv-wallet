package admin

import (
	"errors"
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/mappings"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
// @Deprecated
func transactionRecord(c *gin.Context, _ *reqctx.AdminContext) {
	logger := reqctx.Logger(c)
	engineInstance := reqctx.Engine(c)
	var requestBody RecordTransaction
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest, logger)
		return
	}

	// Set the metadata
	opts := make([]engine.ModelOps, 0)

	// Record a new transaction (get the hex from parameters)
	transaction, err := engineInstance.RecordRawTransaction(
		c.Request.Context(),
		requestBody.Hex,
		opts...,
	)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			// already registered, just return the registered transaction
			if transaction, err = engineInstance.GetTransactionByHex(c.Request.Context(), requestBody.Hex); err != nil {
				spverrors.ErrorResponse(c, err, logger)
				return
			}
		} else {
			spverrors.ErrorResponse(c, err, logger)
			return
		}
	}

	contract := mappings.MapToOldTransactionContract(transaction)

	c.JSON(http.StatusCreated, contract)
}
