package transactions

import (
	"github.com/bitcoin-sv/spv-wallet/actions/v2/transactions/internal/mapping"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	model "github.com/bitcoin-sv/spv-wallet/models/transaction"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func transactionRecordOutline(c *gin.Context, userContext *reqctx.UserContext) {
	logger := reqctx.Logger(c)

	var requestBody model.AnnotatedTransaction
	err := c.ShouldBindWith(&requestBody, binding.JSON)
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest.Wrap(err), logger)
		return
	}

	userID, err := userContext.ShouldGetUserID()
	if err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	}

	recordService := reqctx.Engine(c).TransactionRecordService()
	if err = recordService.RecordTransactionOutline(c, userID, mapping.AnnotatedTransactionToOutline(&requestBody)); err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	}

	c.JSON(200, nil)
}
