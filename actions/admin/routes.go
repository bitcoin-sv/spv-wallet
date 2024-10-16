package admin

import (
	"github.com/bitcoin-sv/spv-wallet/server/handlers"
)

// RegisterRoutes creates the specific package routes
func RegisterRoutes(handlersManager *handlers.Manager) {
	adminGroupOld := handlersManager.Group(handlers.GroupOldAPI, "/admin")

	adminGroupOld.GET("/stats", handlers.AsAdmin(stats))
	adminGroupOld.GET("/status", handlers.AsAdmin(status))
	adminGroupOld.POST("/access-keys/search", handlers.AsAdmin(accessKeysSearch))
	adminGroupOld.POST("/access-keys/count", handlers.AsAdmin(accessKeysCount))
	adminGroupOld.POST("/contact/search", handlers.AsAdmin(contactsSearch))
	adminGroupOld.PATCH("/contact/:id", handlers.AsAdmin(contactsUpdate))
	adminGroupOld.DELETE("/contact/:id", handlers.AsAdmin(contactsDelete))
	adminGroupOld.PATCH("/contact/accepted/:id", handlers.AsAdmin(contactsAccept))
	adminGroupOld.PATCH("/contact/rejected/:id", handlers.AsAdmin(contactsReject))
	adminGroupOld.POST("/destinations/search", handlers.AsAdmin(destinationsSearch))
	adminGroupOld.POST("/destinations/count", handlers.AsAdmin(destinationsCount))
	adminGroupOld.POST("/paymail/get", handlers.AsAdmin(paymailGetAddress))
	adminGroupOld.POST("/paymails/search", handlers.AsAdmin(paymailAddressesSearch))
	adminGroupOld.POST("/paymails/count", handlers.AsAdmin(paymailAddressesCount))
	adminGroupOld.POST("/paymail/create", handlers.AsAdmin(paymailCreateAddress))
	adminGroupOld.DELETE("/paymail/delete", handlers.AsAdmin(paymailDeleteAddress))
	adminGroupOld.POST("/transactions/search", handlers.AsAdmin(transactionsSearch))
	adminGroupOld.POST("/transactions/count", handlers.AsAdmin(transactionsCount))
	adminGroupOld.POST("/transactions/record", handlers.AsAdmin(transactionRecord))
	adminGroupOld.POST("/utxos/search", handlers.AsAdmin(utxosSearch))
	adminGroupOld.POST("/utxos/count", handlers.AsAdmin(utxosCount))
	adminGroupOld.POST("/xpub", handlers.AsAdmin(xpubsCreate))
	adminGroupOld.POST("/xpubs/search", handlers.AsAdmin(xpubsSearch))
	adminGroupOld.POST("/xpubs/count", handlers.AsAdmin(xpubsCount))
	adminGroupOld.POST("/webhooks/subscriptions", handlers.AsAdmin(subscribeWebhook))
	adminGroupOld.DELETE("/webhooks/subscriptions", handlers.AsAdmin(unsubscribeWebhook))
	adminGroupOld.GET("/webhooks/subscriptions", handlers.AsAdmin(getAllWebhooks))
	adminGroupOld.GET("/transactions/:id", handlers.AsAdmin(getTxAdminByIDOld))
	//adminGroupOld.GET("/transactions", handlers.AsAdmin(getTransactionsOld))

	adminGroup := handlersManager.Group(handlers.GroupAPI, "/admin")
	adminGroup.GET("/transactions/:id", handlers.AsAdmin(getTxAdminByID))
	//adminGroup.GET("/transactions", handlers.AsAdmin(getTransactions))
}
