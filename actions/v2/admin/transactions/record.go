package transactions

import (
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"github.com/bsv-blockchain/spv-wallet/actions/v2/admin/internal/mapping"
	"github.com/bsv-blockchain/spv-wallet/api"
	"github.com/bsv-blockchain/spv-wallet/errdef/clienterr"
)

// RecordTransactionOutlineForUser records transaction outline for given user
func (s *APIAdminTransactions) RecordTransactionOutlineForUser(c *gin.Context) {
	var requestBody api.RequestsRecordTransactionOutlineForUser
	err := c.ShouldBindWith(&requestBody, binding.JSON)
	if err != nil {
		clienterr.UnprocessableEntity.New().Wrap(err).Response(c, s.logger)
		return
	}
	if requestBody.UserID == "" {
		clienterr.BadRequest.Detailed("missing_user_id", "userID not provided").Response(c, s.logger)
		return
	}

	outline, err := mapping.RequestsTransactionOutlineToOutline(&requestBody)
	if err != nil {
		clienterr.Response(c, err, s.logger)
		return
	}

	recorded, err := s.transactionsRecordService.RecordTransactionOutline(c.Request.Context(), requestBody.UserID, outline)
	if err != nil {
		clienterr.Response(c, err, s.logger)
		return
	}

	c.JSON(201, mapping.RecordedOutline(recorded))
}
