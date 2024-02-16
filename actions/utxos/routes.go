package utxos

import (
	"github.com/bitcoin-sv/spv-wallet/actions"
	"github.com/bitcoin-sv/spv-wallet/config"
	"github.com/bitcoin-sv/spv-wallet/server/routes"
	"github.com/gin-gonic/gin"
)

// Action is an extension of actions.Action for this package
type Action struct {
	actions.Action
}

// NewHandler creates the specific package routes
func NewHandler(appConfig *config.AppConfig, services *config.AppServices) routes.ApiEndpointsFunc {
	action := &Action{actions.Action{AppConfig: appConfig, Services: services}}

	apiEndpoints := routes.ApiEndpointsFunc(func(router *gin.RouterGroup) {
		utxoGroup := router.Group("/utxo")
		utxoGroup.GET("", action.get)
		utxoGroup.POST("/count", action.count)
		utxoGroup.POST("/search", action.search)
	})

	return apiEndpoints
}
