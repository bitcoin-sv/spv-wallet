package base

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/gin-gonic/gin"
)

// SharedConfig is the handler for SharedConfig which can be obtained by both admin and user
func (s *APIBase) SharedConfig(c *gin.Context) {
	sharedConfig := response.SharedConfig{
		PaymailDomains: s.config.Paymail.Domains,
		ExperimentalFeatures: map[string]bool{
			"pikeContactsEnabled": s.config.ExperimentalFeatures.PikeContactsEnabled,
			"pikePaymentEnabled":  s.config.ExperimentalFeatures.PikePaymentEnabled,
		},
	}

	c.JSON(http.StatusOK, sharedConfig)
}
