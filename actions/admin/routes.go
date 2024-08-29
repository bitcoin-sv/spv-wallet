package admin

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
func NewHandler(appConfig *config.AppConfig, services *config.AppServices) routes.AdminEndpointsFunc {
	action := &Action{actions.Action{AppConfig: appConfig, Services: services}}

	adminEndpoints := routes.AdminEndpointsFunc(func(router *gin.RouterGroup) {
		adminGroup := router.Group("/admin")
		adminGroup.GET("/stats", action.stats)
		adminGroup.GET("/status", action.status)
		adminGroup.POST("/access-keys/search", action.accessKeysSearch)
		adminGroup.POST("/access-keys/count", action.accessKeysCount)
		adminGroup.POST("/contact/search", action.contactsSearch)
		adminGroup.PATCH("/contact/:id", action.contactsUpdate)
		adminGroup.DELETE("/contact/:id", action.contactsDelete)
		adminGroup.PATCH("/contact/accepted/:id", action.contactsAccept)
		adminGroup.PATCH("/contact/rejected/:id", action.contactsReject)
		adminGroup.POST("/destinations/search", action.destinationsSearch)
		adminGroup.POST("/destinations/count", action.destinationsCount)
		adminGroup.POST("/paymail/get", action.paymailGetAddress)
		adminGroup.POST("/paymails/search", action.paymailAddressesSearch)
		adminGroup.POST("/paymails/count", action.paymailAddressesCount)
		adminGroup.POST("/paymail/create", action.paymailCreateAddress)
		adminGroup.DELETE("/paymail/delete", action.paymailDeleteAddress)
		adminGroup.POST("/transactions/search", action.transactionsSearch)
		adminGroup.POST("/transactions/count", action.transactionsCount)
		adminGroup.POST("/transactions/record", action.transactionRecord)
		adminGroup.POST("/utxos/search", action.utxosSearch)
		adminGroup.POST("/utxos/count", action.utxosCount)
		adminGroup.POST("/xpub", action.xpubsCreate)
		adminGroup.POST("/xpubs/search", action.xpubsSearch)
		adminGroup.POST("/xpubs/count", action.xpubsCount)
		adminGroup.POST("/webhooks/subscriptions", action.subscribeWebhook)
		adminGroup.DELETE("/webhooks/subscriptions", action.unsubscribeWebhook)
		adminGroup.GET("/webhooks/subscriptions", action.getAllWebhooks)
	})

	return adminEndpoints
}
