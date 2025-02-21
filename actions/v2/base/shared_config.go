package base

import (
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors/examplecode/domain/exampledomain"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors/examplecode/errdef"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors/examplecode/repos"
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/api"
	"github.com/gin-gonic/gin"
)

// SharedConfig is the handler for SharedConfig which can be obtained by both admin and user
func (s *APIBase) SharedConfig(c *gin.Context) {
	sharedConfig := api.ResponsesSharedConfig{
		PaymailDomains: s.config.Paymail.Domains,
		ExperimentalFeatures: map[string]bool{
			"pikeContactsEnabled": s.config.ExperimentalFeatures.PikeContactsEnabled,
			"pikePaymentEnabled":  s.config.ExperimentalFeatures.PikePaymentEnabled,
			"v2":                  s.config.ExperimentalFeatures.V2,
		},
	}

	c.JSON(http.StatusOK, sharedConfig)
}

func (s *APIBase) GetApiV2TestErrors(c *gin.Context, params api.GetApiV2TestErrorsParams) {

	repo := repos.NewRepo()
	domain := exampledomain.NewService(repo)

	_, err := domain.Search(params.Fail)
	if err != nil {
		problemDetails := errdef.NewProblemDetailsFromError(err)
		c.JSON(problemDetails.Status, problemDetails)
		return
	}

	c.JSON(200, "success") // not relevant for the example
}
