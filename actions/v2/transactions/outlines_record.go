package transactions

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"github.com/bsv-blockchain/spv-wallet/actions/v2/transactions/internal/mapping"
	"github.com/bsv-blockchain/spv-wallet/api"
	"github.com/bsv-blockchain/spv-wallet/engine/spverrors"
	"github.com/bsv-blockchain/spv-wallet/server/reqctx"
)

// RecordTransactionOutline records transaction outline
func (s *APITransactions) RecordTransactionOutline(c *gin.Context) {
	var requestBody api.RequestsTransactionOutline
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

	outline, err := mapping.RequestsTransactionOutlineToOutline(&requestBody)
	if err != nil {
		spverrors.ErrorResponse(c, err, s.logger)
		return
	}

	recordService := s.engine.TransactionRecordService()
	recorded, err := recordService.RecordTransactionOutline(c.Request.Context(), userID, outline)
	if err != nil {
		spverrors.ErrorResponse(c, err, s.logger)
		return
	}

	c.JSON(201, mapping.RecordedOutline(recorded))
}
