package base

import (
	"github.com/bitcoin-sv/spv-wallet/actions"
	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/bitcoin-sv/spv-wallet/server/routes"
	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
)

// Action is an extension of actions.Action for this package
type Action struct {
	actions.Action
}

// NewHandler creates the specific package routes
func NewHandler(appConfig *config.AppConfig, engine *gin.Engine) routes.BaseEndpointsFunc {
	basicEndpoints := routes.BaseEndpointsFunc(func(router *gin.RouterGroup) {
		router.GET("/", index)
		router.OPTIONS("/", actions.StatusOK)
		router.HEAD("/", actions.StatusOK)

		healthGroup := router.Group("/" + config.HealthRequestPath)
		healthGroup.GET("", actions.StatusOK)
		healthGroup.OPTIONS("", actions.StatusOK)
		healthGroup.HEAD("", actions.StatusOK)
	})

	if appConfig.DebugProfiling {
		pprof.Register(engine, "debug/pprof")
	}

	return basicEndpoints
}
