package transactions

import (
	"github.com/bitcoin-sv/spv-wallet/actions/v2/transactions/internal/mapping"
	"github.com/bitcoin-sv/spv-wallet/api"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// RecordTransactionOutline records transaction outline
func (s *APITransactions) RecordTransactionOutline(c *gin.Context) {
	var requestBody api.RequestsAnnotatedTransaction
	err := c.ShouldBindWith(&requestBody, binding.JSON)
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest.Wrap(err), s.logger)
		return
	}

	userContext := reqctx.GetUserContext(c)
	userID, err := userContext.ShouldGetUserID()
	if err != nil {
		spverrors.ErrorResponse(c, err, s.logger)
		return
	}

	outline, err := mapping.AnnotatedTransactionRequestToOutline(&requestBody)
	if err != nil {
		spverrors.ErrorResponse(c, err, s.logger)
		return
	}

	recordService := s.engine.TransactionRecordService()
	recorded, err := recordService.RecordTransactionOutline(c, userID, outline)
	if err != nil {
		spverrors.ErrorResponse(c, err, s.logger)
		return
	}

	c.JSON(201, mapping.RecordedOutline(recorded))
}
