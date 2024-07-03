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

// NewTransactionsHandler creates the specific package routes in restfull style
func NewTransactionsHandler(appConfig *config.AppConfig, services *config.AppServices) *routes.Handler {
	action := &Action{actions.Action{AppConfig: appConfig, Services: services}}

	handler := routes.Handler{
		BasicEndpointsFunc: routes.BasicEndpointsFunc(func(router *gin.RouterGroup) {
			basicTransactionGroup := router.Group("/transactions")
			basicTransactionGroup.GET(":id", action.getByID)
			basicTransactionGroup.PATCH(":id", action.updateTransaction)
			basicTransactionGroup.GET("", action.transactions)
		}),
		APIEndpoints: routes.APIEndpointsFunc(func(router *gin.RouterGroup) {
			apiTransactionGroup := router.Group("/transactions")
			apiTransactionGroup.POST("/drafts", action.newTransactionDraft)
			apiTransactionGroup.POST("", action.recordTransaction)
		}),
		CallbackEndpoints: routes.CallbackEndpointsFunc(func(router *gin.RouterGroup) {
			router.POST(config.BroadcastCallbackRoute, action.broadcastCallback)
		}),
	}
	return &handler
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
