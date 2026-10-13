package webhooks

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/bsv-blockchain/spv-wallet/actions/v2/admin/internal/mapping"
	"github.com/bsv-blockchain/spv-wallet/engine/spverrors"
)

// Webhooks returns all webhooks
func (s *APIAdminWebhooks) Webhooks(c *gin.Context) {
	wh, err := s.webhooks.GetWebhooks(c.Request.Context())
	if err != nil {
		spverrors.ErrorResponse(c, err, s.logger)
		return
	}

	c.JSON(http.StatusOK, mapping.MapToModelsWebhooks(wh))
}
