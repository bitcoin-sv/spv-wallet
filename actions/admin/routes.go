package admin

import (
	"github.com/bitcoin-sv/spv-wallet/server/handlers"
)

// RegisterRoutes creates the specific package routes
func RegisterRoutes(handlersManager *handlers.Manager) {
	adminGroup := handlersManager.Group(handlers.GroupOldAPI, "/admin")

	adminGroup.GET("/stats", handlers.AsAdmin(stats))
	adminGroup.GET("/status", handlers.AsAdmin(status))
	adminGroup.POST("/access-keys/search", handlers.AsAdmin(accessKeysSearch))
	adminGroup.POST("/access-keys/count", handlers.AsAdmin(accessKeysCount))
	adminGroup.POST("/contact/search", handlers.AsAdmin(contactsSearch))
	adminGroup.PATCH("/contact/:id", handlers.AsAdmin(contactsUpdate))
	adminGroup.DELETE("/contact/:id", handlers.AsAdmin(contactsDelete))
	adminGroup.PATCH("/contact/accepted/:id", handlers.AsAdmin(contactsAccept))
	adminGroup.PATCH("/contact/rejected/:id", handlers.AsAdmin(contactsReject))
	adminGroup.POST("/destinations/search", handlers.AsAdmin(destinationsSearch))
	adminGroup.POST("/destinations/count", handlers.AsAdmin(destinationsCount))
	adminGroup.POST("/paymail/get", handlers.AsAdmin(paymailGetAddress))
	adminGroup.POST("/paymails/search", handlers.AsAdmin(paymailAddressesSearch))
	adminGroup.POST("/paymails/count", handlers.AsAdmin(paymailAddressesCount))
	adminGroup.POST("/paymail/create", handlers.AsAdmin(paymailCreateAddress))
	adminGroup.DELETE("/paymail/delete", handlers.AsAdmin(paymailDeleteAddress))
	adminGroup.POST("/transactions/search", handlers.AsAdmin(transactionsSearch))
	adminGroup.POST("/transactions/count", handlers.AsAdmin(transactionsCount))
	adminGroup.POST("/transactions/record", handlers.AsAdmin(transactionRecord))
	adminGroup.POST("/utxos/search", handlers.AsAdmin(utxosSearch))
	adminGroup.POST("/utxos/count", handlers.AsAdmin(utxosCount))
	adminGroup.POST("/xpub", handlers.AsAdmin(xpubsCreate))
	adminGroup.POST("/xpubs/search", handlers.AsAdmin(xpubsSearch))
	adminGroup.POST("/xpubs/count", handlers.AsAdmin(xpubsCount))
	adminGroup.POST("/webhooks/subscriptions", handlers.AsAdmin(subscribeWebhook))
	adminGroup.DELETE("/webhooks/subscriptions", handlers.AsAdmin(unsubscribeWebhook))
	adminGroup.GET("/webhooks/subscriptions", handlers.AsAdmin(getAllWebhooks))
}
