package admin

import (
	"github.com/bitcoin-sv/spv-wallet/server/handlers"
)

// RegisterRoutes creates the specific package routes
func RegisterRoutes(handlersManager *handlers.Manager) {
	adminGroupOld := handlersManager.Group(handlers.GroupOldAPI, "/admin")

	adminGroupOld.POST("/xpub", handlers.AsAdmin(xpubsCreateOld))
	adminGroupOld.POST("/xpubs/search", handlers.AsAdmin(xpubsSearchOld))
	adminGroupOld.POST("/xpubs/count", handlers.AsAdmin(xpubsCount))

	adminGroupOld.POST("/webhooks/subscriptions", handlers.AsAdmin(subscribeWebhookOld))
	adminGroupOld.DELETE("/webhooks/subscriptions", handlers.AsAdmin(unsubscribeWebhookOld))
	adminGroupOld.GET("/webhooks/subscriptions", handlers.AsAdmin(getAllWebhooksOld))

	adminGroup := handlersManager.Group(handlers.GroupAPI, "/admin")
	adminGroup.GET("/status", handlers.AsAdmin(status))
	adminGroup.GET("/stats", handlers.AsAdmin(stats))

	// tx
	adminGroup.GET("/transactions/:id", handlers.AsAdmin(adminGetTxByID))
	adminGroup.GET("/transactions", handlers.AsAdmin(adminSearchTxs))

	// contacts
	adminGroup.GET("/contacts", handlers.AsAdmin(contactsSearch))
	adminGroup.POST("/invitations/:id", handlers.AsAdmin(contactsAccept))
	adminGroup.DELETE("/invitations/:id", handlers.AsAdmin(contactsReject))
	adminGroup.DELETE("/contacts/:id", handlers.AsAdmin(contactsDelete))
	adminGroup.PUT("/contacts/:id", handlers.AsAdmin(contactsUpdate))
	adminGroup.POST("/contacts/:paymail", handlers.AsAdmin(contactsCreate))
	adminGroup.POST("/contacts/confirmations", handlers.AsAdmin(contactsConfirm))

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
	adminGroup.GET("/webhooks/subscriptions", handlers.AsAdmin(getAllWebhooks))
	adminGroup.POST("/webhooks/subscriptions", handlers.AsAdmin(subscribeWebhook))
	adminGroup.DELETE("/webhooks/subscriptions", handlers.AsAdmin(unsubscribeWebhook))

	// xpubs => users
	adminGroup.POST("/users", handlers.AsAdmin(xpubsCreate)) // create
	adminGroup.GET("/users", handlers.AsAdmin(xpubsSearch))  // search
}
