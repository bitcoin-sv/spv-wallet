package sharedconfig

import (
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
	"net/http"
)

// sharedConfig will return the shared configuration
// Get shared config godoc
// @Summary		Get shared config
// @Description	Get shared config
// @Tags		Configurations
// @Produce		json
// @Success		200 {object} response.SharedConfig "Shared configuration"
// @Router		/api/v1/configs/shared [get]
// @Security	x-auth-xpub
func get(c *gin.Context, _ *reqctx.UserContext) {
	appConfig := reqctx.AppConfig(c)
	sharedConfig := response.SharedConfig{
		PaymailDomains: appConfig.Paymail.Domains,
		ExperimentalFeatures: map[string]bool{
			"pikeContactsEnabled": appConfig.ExperimentalFeatures.PikeContactsEnabled,
			"pikePaymentEnabled":  appConfig.ExperimentalFeatures.PikePaymentEnabled,
		},
	}

	c.JSON(http.StatusOK, sharedConfig)
}
