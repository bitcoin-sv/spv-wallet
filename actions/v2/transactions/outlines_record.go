package transactions

import (
	"github.com/bitcoin-sv/spv-wallet/actions/v2/transactions/internal/mapping"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/models/request"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func recordOutline(c *gin.Context, userContext *reqctx.UserContext) {
	logger := reqctx.Logger(c)

	var requestBody request.AnnotatedTransaction
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
	recorded, err := recordService.RecordTransactionOutline(c, userID, mapping.TransactionOutline(&requestBody))
	if err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	}

	c.JSON(201, mapping.RecordedOutline(recorded))
}
