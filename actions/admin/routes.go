package admin

import (
	"github.com/bitcoin-sv/spv-wallet/server/handlers"
)

// RegisterRoutes creates the specific package routes
func RegisterRoutes(handlersManager *handlers.Manager) {
	adminGroupOld := handlersManager.Group(handlers.GroupOldAPI, "/admin")

	adminGroupOld.GET("/stats", handlers.AsAdmin(stats))
	adminGroupOld.GET("/status", handlers.AsAdmin(status))
	adminGroupOld.POST("/access-keys/search", handlers.AsAdmin(accessKeysSearchOld))
	adminGroupOld.POST("/access-keys/count", handlers.AsAdmin(accessKeysCount))
	adminGroupOld.POST("/contact/search", handlers.AsAdmin(contactsSearchOld))
	adminGroupOld.PATCH("/contact/:id", handlers.AsAdmin(contactsUpdateOld))
	adminGroupOld.DELETE("/contact/:id", handlers.AsAdmin(contactsDeleteOld))
	adminGroupOld.PATCH("/contact/accepted/:id", handlers.AsAdmin(contactsAcceptOld))
	adminGroupOld.PATCH("/contact/rejected/:id", handlers.AsAdmin(contactsRejectOld))
	adminGroupOld.POST("/destinations/search", handlers.AsAdmin(destinationsSearch))
	adminGroupOld.POST("/destinations/count", handlers.AsAdmin(destinationsCount))
	adminGroupOld.POST("/paymail/get", handlers.AsAdmin(paymailGetAddressOld))
	adminGroupOld.POST("/paymails/search", handlers.AsAdmin(paymailAddressesSearchOld))
	adminGroupOld.POST("/paymails/count", handlers.AsAdmin(paymailAddressesCount))
	adminGroupOld.POST("/paymail/create", handlers.AsAdmin(paymailCreateAddressOld))
	adminGroupOld.DELETE("/paymail/delete", handlers.AsAdmin(paymailDeleteAddressOld))
	adminGroupOld.POST("/transactions/search", handlers.AsAdmin(transactionsSearch))
	adminGroupOld.POST("/transactions/count", handlers.AsAdmin(transactionsCount))
	adminGroupOld.POST("/transactions/record", handlers.AsAdmin(transactionRecord))
	adminGroupOld.POST("/utxos/search", handlers.AsAdmin(utxosSearchOld))
	adminGroupOld.POST("/utxos/count", handlers.AsAdmin(utxosCount))
	adminGroupOld.POST("/xpub", handlers.AsAdmin(xpubsCreateOld))
	adminGroupOld.POST("/xpubs/search", handlers.AsAdmin(xpubsSearchOld))
	adminGroupOld.POST("/xpubs/count", handlers.AsAdmin(xpubsCount))
	adminGroupOld.POST("/webhooks/subscriptions", handlers.AsAdmin(subscribeWebhookOld))
	adminGroupOld.DELETE("/webhooks/subscriptions", handlers.AsAdmin(unsubscribeWebhookOld))
	adminGroupOld.GET("/webhooks/subscriptions", handlers.AsAdmin(getAllWebhooks))

	adminGroupOld.GET("/transactions/:id", handlers.AsAdmin(getTxAdminByIDOld))
	adminGroupOld.GET("/transactions", handlers.AsAdmin(getTransactionsOld))

	adminGroup := handlersManager.Group(handlers.GroupAPI, "/admin")
	// tx
	adminGroup.GET("/transactions/:id", handlers.AsAdmin(adminGetTxByID))
	adminGroup.GET("/transactions", handlers.AsAdmin(adminSearchTxs))

	// contacts
	adminGroup.GET("/contacts", handlers.AsAdmin(contactsSearch))
	adminGroup.POST("/invitations/:id", handlers.AsAdmin(contactsAccept))
	adminGroup.DELETE("/invitations/:id", handlers.AsAdmin(contactsReject))
	adminGroup.DELETE("/contacts/:id", handlers.AsAdmin(contactsDelete))
	adminGroup.PUT("/contacts/:id", handlers.AsAdmin(contactsUpdate))

	// access keys
	adminGroup.GET("/users/keys", handlers.AsAdmin(accessKeysSearch))

	// paymails
	adminGroup.GET("/paymails/:id", handlers.AsAdmin(paymailGetAddress))
	adminGroup.GET("/paymails", handlers.AsAdmin(paymailAddressesSearch))
	adminGroup.POST("/paymails", handlers.AsAdmin(paymailCreateAddress))
	adminGroup.DELETE("/paymails/:id", handlers.AsAdmin(paymailDeleteAddress))

	// utxos
	adminGroup.GET("/utxos", handlers.AsAdmin(utxosSearch))

	// webhooks
	adminGroup.POST("/webhooks/subscriptions", handlers.AsAdmin(subscribeWebhook))
	adminGroup.DELETE("/webhooks/subscriptions", handlers.AsAdmin(unsubscribeWebhook))

	// xpubs => users
	adminGroup.POST("/users", handlers.AsAdmin(xpubsCreate)) // create
	adminGroup.GET("/users", handlers.AsAdmin(xpubsSearch))  // search
}
