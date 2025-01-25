package transactions

import (
	"github.com/bitcoin-sv/spv-wallet/actions/transactions/internal/mapping/outline"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/models/request"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func transactionOutlines(c *gin.Context, userCtx *reqctx.UserContext) {
	logger := reqctx.Logger(c)

	var requestBody request.TransactionSpecification
	err := c.ShouldBindWith(&requestBody, binding.JSON)
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest.Wrap(err), logger)
		return
	}

	userID, err := userCtx.ShouldGetUserID()
	if err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	}

	spec, err := outline.Request(requestBody).ToEngine(userID)
	if err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	}

	txOutline, err := reqctx.Engine(c).TransactionOutlinesService().Create(c, spec)
	if err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	}

	res, err := outline.ToResponse(txOutline)
	if err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	}
	c.JSON(200, res)
}
