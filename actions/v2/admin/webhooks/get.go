package webhooks

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/actions/v2/admin/internal/mapping"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/gin-gonic/gin"
)

// Webhooks returns all webhooks
func (s *APIAdminWebhooks) Webhooks(c *gin.Context) {
	wh, err := s.webhooks.GetWebhooks(c)
	if err != nil {
		spverrors.ErrorResponse(c, err, s.logger)
		return
	}

	c.JSON(http.StatusOK, mapping.MapToModelsWebhooks(wh))
}
