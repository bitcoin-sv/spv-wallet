package contacts

import (
	"context"

	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/contacts/contactsmodels"
	"github.com/rs/zerolog"
)

type contactsService interface {
	AcceptContactByID(ctx context.Context, id uint) (*contactsmodels.Contact, error)
	AdminConfirmContacts(ctx context.Context, paymailA string, paymailB string) error
	AdminCreateContact(ctx context.Context, newContact contactsmodels.NewContact) (*contactsmodels.Contact, error)
	RejectContactByID(ctx context.Context, id uint) (*contactsmodels.Contact, error)
	RemoveContactByID(ctx context.Context, contactID uint) error
	UpdateFullNameByID(ctx context.Context, contactID uint, fullName string) (*contactsmodels.Contact, error)
}

// APIAdminContacts represents server with admin API endpoints
type APIAdminContacts struct {
	contactsService contactsService
	logger          *zerolog.Logger
}

// NewAPIAdminContacts creates a new APIAdminUsers
func NewAPIAdminContacts(engine engine.ClientInterface, logger *zerolog.Logger) APIAdminContacts {
	return APIAdminContacts{
		contactsService: engine.ContactService(),
		logger:          logger,
	}
}
