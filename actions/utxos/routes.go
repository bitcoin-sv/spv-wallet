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
func OldUtxosHandler(appConfig *config.AppConfig, services *config.AppServices) routes.OldAPIEndpointsFunc {
	action := &Action{actions.Action{AppConfig: appConfig, Services: services}}

	apiEndpoints := routes.OldAPIEndpointsFunc(func(router *gin.RouterGroup) {
		utxoGroup := router.Group("/utxo")
		utxoGroup.GET("", action.get)
		utxoGroup.POST("/count", action.count)
		utxoGroup.POST("/search", action.oldSearch)
	})

	return apiEndpoints
}

// NewHandler creates the specific package routes
func NewHandler(appConfig *config.AppConfig, services *config.AppServices) routes.APIEndpointsFunc {
	action := &Action{actions.Action{AppConfig: appConfig, Services: services}}

	apiEndpoints := routes.APIEndpointsFunc(func(router *gin.RouterGroup) {
		utxoGroup := router.Group("/utxos")
		utxoGroup.GET("", action.get)
		utxoGroup.POST("/search", action.search)
	})

	return apiEndpoints
}
