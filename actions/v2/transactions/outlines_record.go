package transactions

import (
	"github.com/bitcoin-sv/spv-wallet/actions/v2/transactions/internal/mapping"
	"github.com/bitcoin-sv/spv-wallet/api"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func (s *APITransactions) PostApiV2Transactions(c *gin.Context) {
	logger := reqctx.Logger(c)

	var requestBody api.ApiComponentsRequestsAnnotatedTransaction
	err := c.ShouldBindWith(&requestBody, binding.JSON)
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest.Wrap(err), logger)
		return
	}

	userContext := reqctx.GetUserContext(c)
	userID, err := userContext.ShouldGetUserID()
	if err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	}

	recordService := reqctx.Engine(c).TransactionRecordService()
	recorded, err := recordService.RecordTransactionOutline(c, userID, mapping.AnnotatedTransactionRequestToOutline(&requestBody))
	if err != nil {
		spverrors.ErrorResponse(c, err, logger)
		return
	}

	c.JSON(201, mapping.RecordedOutline(recorded))
}
