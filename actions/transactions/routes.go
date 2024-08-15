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

// NewHandler creates the specific package routes in restfull style
func NewHandler(appConfig *config.AppConfig, services *config.AppServices) *routes.Handler {
	action := &Action{actions.Action{AppConfig: appConfig, Services: services}}

	handler := routes.Handler{
		BasicEndpoints: routes.BasicEndpointsFunc(func(router *gin.RouterGroup) {
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

// OldTransactionsHandler creates the specific package routes
func OldTransactionsHandler(appConfig *config.AppConfig, services *config.AppServices) (routes.OldBasicEndpointsFunc, routes.OldAPIEndpointsFunc, routes.CallbackEndpointsFunc) {
	action := &Action{actions.Action{AppConfig: appConfig, Services: services}}

	basicEndpoints := routes.OldBasicEndpointsFunc(func(router *gin.RouterGroup) {
		basicTransactionGroup := router.Group("/transaction")
		basicTransactionGroup.GET("", action.get)
		basicTransactionGroup.PATCH("", action.update)
		basicTransactionGroup.POST("/count", action.count)
		basicTransactionGroup.GET("/search", action.search)
		basicTransactionGroup.POST("/search", action.search)
	})

	apiEndpoints := routes.OldAPIEndpointsFunc(func(router *gin.RouterGroup) {
		apiTransactionGroup := router.Group("/transaction")
		apiTransactionGroup.POST("", action.newTransaction)
		apiTransactionGroup.POST("/record", action.record)
	})

	callbackEndpoints := routes.CallbackEndpointsFunc(func(router *gin.RouterGroup) {
		router.POST(config.BroadcastCallbackRoute, action.broadcastCallback)
	})

	return basicEndpoints, apiEndpoints, callbackEndpoints
}
