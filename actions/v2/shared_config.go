package v2

import (
	"sync"

	"github.com/bitcoin-sv/spv-wallet/api"
	"github.com/bitcoin-sv/spv-wallet/server/reqctx"
	"github.com/gin-gonic/gin"
)

// GetApiV2ConfigsShared is the handler for SharedConfig which can be obtained by both admin and user
func (s *Server) GetApiV2ConfigsShared(c *gin.Context) {
	appconfig := reqctx.AppConfig(c)
	makeConfig := sync.OnceValue(func() api.ApiComponentsResponsesSharedConfig {
		return api.ApiComponentsResponsesSharedConfig{
			PaymailDomains: &appconfig.Paymail.Domains,
			ExperimentalFeatures: &map[string]bool{
				"pikeContactsEnabled": appconfig.ExperimentalFeatures.PikeContactsEnabled,
				"pikePaymentEnabled":  appconfig.ExperimentalFeatures.PikePaymentEnabled,
			},
		}
	})

	c.JSON(200, makeConfig())
}
