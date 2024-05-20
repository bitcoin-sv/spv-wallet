package sharedconfig

import (
	"net/http"
	"sync"

	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/gin-gonic/gin"
)

// sharedConfig will return the shared configuration
// Get shared config godoc
// @Summary		Get shared config
// @Description	Get shared config
// @Tags		Shared-config
// @Produce		json
// @Success		200 {object} models.SharedConfig "Shared configuration"
// @Router		/v1/shared-config [get]
// @Security	x-auth-xpub
func (a *Action) get(c *gin.Context) {
	makeConfig := sync.OnceValue(func() models.SharedConfig {
		return models.SharedConfig{
			PaymailDomains: a.AppConfig.Paymail.Domains,
			ExperimentalFeatures: map[string]bool{
				"pike_contacts_enabled": a.AppConfig.ExperimentalFeatures.PikeContactsEnabled,
				"pike_payment_enabled":  a.AppConfig.ExperimentalFeatures.PikePaymentEnabled,
			},
		}
	})

	c.JSON(http.StatusOK, makeConfig())
}
