package admin

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/gin-gonic/gin"
)

// subscribeWebhook will subscribe to a webhook to receive notifications
// @Summary		Subscribe to a webhook
// @Description	Subscribe to a webhook to receive notifications
// @Tags		Admin
// @Produce		json
// @Param		SubscribeRequestBody body models.SubscribeRequestBody false "URL to subscribe to and optional token header and value"
// @Success		200 {boolean} bool "Success response"
// @Failure 	500	"Internal server error - Error while subscribing to the webhook"
// @Router		/v1/admin/webhooks/subscription [post]
// @Security	x-auth-xpub
func (a *Action) subscribeWebhook(c *gin.Context) {
	requestBody := models.SubscribeRequestBody{}
	if err := c.Bind(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err := a.Services.SpvWalletEngine.SubscribeWebhook(c.Request.Context(), requestBody.URL, requestBody.TokenHeader, requestBody.TokenValue)
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrWebhookSubscriptionFailed, a.Services.Logger)
		return
	}

	c.JSON(http.StatusOK, true)
}

// unsubscribeWebhook will unsubscribe to a webhook to receive notifications
// @Summary		Unsubscribe to a webhook
// @Description	Unsubscribe to a webhook to stop receiving notifications
// @Tags		Admin
// @Produce		json
// @Param		UnsubscribeRequestBody body models.UnsubscribeRequestBody false "URL to unsubscribe from"
// @Success		200 {boolean} bool "Success response"
// @Failure 	500	"Internal server error - Error while unsubscribing to the webhook"
// @Router		/v1/admin/webhooks/subscription [delete]
// @Security	x-auth-xpub
func (a *Action) unsubscribeWebhook(c *gin.Context) {
	requestModel := models.UnsubscribeRequestBody{}
	if err := c.Bind(&requestModel); err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}

	err := a.Services.SpvWalletEngine.UnsubscribeWebhook(c.Request.Context(), requestModel.URL)
	if err != nil {
		spverrors.ErrorResponse(c, spverrors.ErrWebhookUnsubscriptionFailed, a.Services.Logger)
		return
	}

	c.JSON(http.StatusOK, true)
}
