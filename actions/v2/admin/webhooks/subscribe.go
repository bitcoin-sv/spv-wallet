package webhooks

import (
	"net/http"
	"net/url"

	"github.com/bitcoin-sv/spv-wallet/api"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// SubscribeWebhook subscribes to a webhook
func (s *APIAdminWebhooks) SubscribeWebhook(c *gin.Context) {
	var bodyReq api.RequestsSubscribeWebhook
	if err := c.ShouldBindWith(&bodyReq, binding.JSON); err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrCannotBindRequest.WithTrace(err), s.logger)
		return
	}

	if bodyReq.Url == "" {
		spverrors.ErrorResponse(c, spverrors.ErrWebhookUrlRequired, s.logger)
		return
	}

	if _, err := url.Parse(bodyReq.Url); err != nil {
		spverrors.ErrorResponse(c, spverrors.WebhookUrlInvalid, s.logger)
		return
	}

	if bodyReq.TokenHeader == "" {
		spverrors.ErrorResponse(c, spverrors.ErrWebhookTokenHeaderRequired, s.logger)
		return
	}

	if bodyReq.TokenValue == "" {
		spverrors.ErrorResponse(c, spverrors.ErrWebhookTokenValueRequired, s.logger)
		return
	}

	err := s.webhooks.SubscribeWebhook(c, bodyReq.Url, bodyReq.TokenHeader, bodyReq.TokenValue)
	if err != nil {
		spverrors.ErrorResponse(c, err, s.logger)
		return
	}

	c.Status(http.StatusOK)
}
