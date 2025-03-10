package transactions

import (
	"github.com/bitcoin-sv/spv-wallet/actions/v2/admin/internal/mapping"
	"github.com/bitcoin-sv/spv-wallet/api"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// RecordTransactionOutlineForUser records transaction outline for given user
func (s *APIAdminTransactions) RecordTransactionOutlineForUser(c *gin.Context) {
	var requestBody api.RequestsRecordTransactionOutlineForUser
	err := c.ShouldBindWith(&requestBody, binding.JSON)
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest.Wrap(err), s.logger)
		return
	}
	if requestBody.UserID == "" {
		spverrors.ErrorResponse(c, spverrors.Wrapf(spverrors.ErrCannotBindRequest, "userID not provided"), s.logger)
		return
	}

	outline, err := mapping.RequestsTransactionOutlineToOutline(&requestBody)
	if err != nil {
		spverrors.ErrorResponse(c, err, s.logger)
		return
	}

	recorded, err := s.transactionsRecordService.RecordTransactionOutline(c, requestBody.UserID, outline)
	if err != nil {
		spverrors.ErrorResponse(c, err, s.logger)
		return
	}

	c.JSON(201, mapping.RecordedOutline(recorded))
}
