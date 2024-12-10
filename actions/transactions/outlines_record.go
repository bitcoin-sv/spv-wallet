package transactions

import (
	"github.com/bitcoin-sv/spv-wallet/actions/transactions/internal/mapping/annotatedtx"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func transactionRecordOutline(c *gin.Context, userContext *reqctx.UserContext) {
	logger := reqctx.Logger(c)

	var requestBody annotatedtx.Request
	err := c.ShouldBindWith(&requestBody, binding.JSON)
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest.Wrap(err), logger)
		return
	}

	recordService := reqctx.Engine(c).TransactionRecordService()
	if err = recordService.RecordTransactionOutline(c, userContext.GetXPubID(), requestBody.ToEngine()); err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	}

	c.JSON(200, nil)
}
