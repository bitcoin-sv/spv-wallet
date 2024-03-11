package admin

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
// @Tags		Admin
// @Produce		json
// @Success		200
// @Router		/v1/admin/shared-config [get]
// @Security	x-auth-xpub
func (a *Action) sharedConfig(c *gin.Context) {
	makeConfig := sync.OnceValue(func() models.SharedConfig {
		return models.SharedConfig{
			PaymilDomains: a.AppConfig.Paymail.Domains,
			ExperimentalFeatures: map[string]bool{
				"pike_enabled": a.AppConfig.ExperimentalFeatures.PikeEnabled,
			},
		}
	})

	c.JSON(http.StatusOK, makeConfig())
}
