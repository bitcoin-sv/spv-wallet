package base

import (
	"github.com/bitcoin-sv/spv-wallet/models/response"
	"github.com/gin-gonic/gin"
	"net/http"
)

// GetApiV2ConfigsShared is the handler for SharedConfig which can be obtained by both admin and user
func (s *APIBase) GetApiV2ConfigsShared(c *gin.Context) {
	sharedConfig := response.SharedConfig{
		PaymailDomains: s.config.Paymail.Domains,
		ExperimentalFeatures: map[string]bool{
			"pikeContactsEnabled": s.config.ExperimentalFeatures.PikeContactsEnabled,
			"pikePaymentEnabled":  s.config.ExperimentalFeatures.PikePaymentEnabled,
		},
	}

	c.JSON(http.StatusOK, sharedConfig)
}
