package transactions

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
func NewHandler(appConfig *config.AppConfig, services *config.AppServices) (routes.BasicEndpointsFunc, routes.APIEndpointsFunc, routes.CallbackEndpointsFunc) {
	action := &Action{actions.Action{AppConfig: appConfig, Services: services}}

	basicEndpoints := routes.BasicEndpointsFunc(func(router *gin.RouterGroup) {
		basicTransactionGroup := router.Group("/transaction")
		basicTransactionGroup.GET("", action.get)
		basicTransactionGroup.PATCH("", action.update)
		basicTransactionGroup.POST("/count", action.count)
		basicTransactionGroup.GET("/search", action.search)
		basicTransactionGroup.POST("/search", action.search)
	})

	apiEndpoints := routes.APIEndpointsFunc(func(router *gin.RouterGroup) {
		apiTransactionGroup := router.Group("/transaction")
		apiTransactionGroup.POST("", action.newTransaction)
		apiTransactionGroup.POST("/record", action.record)
	})

	callbackEndpoints := routes.CallbackEndpointsFunc(func(router *gin.RouterGroup) {
		router.POST(config.BroadcastCallbackRoute, action.broadcastCallback)
	})

	return basicEndpoints, apiEndpoints, callbackEndpoints
}
