package contacts

import (
	"github.com/bitcoin-sv/spv-wallet/server/handlers"
)

// RegisterRoutes creates the specific package routes
func RegisterRoutes(handlersManager *handlers.Manager) {
	groupContacts := handlersManager.Group(handlers.GroupAPI, "/contacts")
	groupContacts.PUT("/:paymail", handlers.AsUser(upsertContact))
	groupContacts.DELETE("/:paymail", handlers.AsUser(removeContact))

	groupContacts.POST("/:paymail/confirmation", handlers.AsUser(confirmContact))
	groupContacts.DELETE("/:paymail/confirmation", handlers.AsUser(unconfirmContact))

	groupContacts.GET("", handlers.AsUser(getContacts))
	groupContacts.GET(":paymail", handlers.AsUser(getContactByPaymail))

	groupInvitations := handlersManager.Group(handlers.GroupAPI, "/invitations")
	groupInvitations.POST("/:paymail/contacts", handlers.AsUser(acceptInvitations))
	groupInvitations.DELETE("/:paymail", handlers.AsUser(rejectInvitation))
}
