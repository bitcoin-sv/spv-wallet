package contacts

import (
	"github.com/bitcoin-sv/spv-wallet/server/handlers"
)

// RegisterRoutes creates the specific package routes
func RegisterRoutes(handlersManager *handlers.Manager) {
	old := handlersManager.Group(handlers.GroupOldAPI, "/contact")
	old.PUT("/:paymail", handlers.AsUser(oldUpsert))
	old.PATCH("/accepted/:paymail", handlers.AsUser(oldAccept))
	old.PATCH("/rejected/:paymail", handlers.AsUser(oldReject))
	old.PATCH("/confirmed/:paymail", handlers.AsUser(oldConfirm))
	old.PATCH("/unconfirmed/:paymail", handlers.AsUser(oldUnconfirm))
	old.POST("search", handlers.AsUser(search))

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
