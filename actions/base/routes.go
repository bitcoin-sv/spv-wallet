package base

import (
	"github.com/bitcoin-sv/spv-wallet/actions"
	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/bitcoin-sv/spv-wallet/server/routes"
	"github.com/gin-gonic/gin"
)

// NewHandler creates the specific package routes
func NewHandler() routes.BaseEndpointsFunc {
	basicEndpoints := routes.BaseEndpointsFunc(func(router *gin.RouterGroup) {
		router.GET("/", index)
		router.OPTIONS("/", actions.StatusOK)
		router.HEAD("/", actions.StatusOK)

		healthGroup := router.Group("/" + config.HealthRequestPath)
		healthGroup.GET("", actions.StatusOK)
		healthGroup.OPTIONS("", actions.StatusOK)
		healthGroup.HEAD("", actions.StatusOK)
	})

	return basicEndpoints
}
