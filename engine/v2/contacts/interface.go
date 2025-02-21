package contacts

import (
	"context"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/contacts/contactsmodels"
)

// ContactRepo is a contacts repository
type ContactRepo interface {
	Create(ctx context.Context, newContact *contactsmodels.NewContact) (*contactsmodels.Contact, error)
	Update(ctx context.Context, contact *contactsmodels.NewContact) (*contactsmodels.Contact, error)
	Delete(ctx context.Context, userID, paymail string) error
	Find(ctx context.Context, userID, paymail string) (*contactsmodels.Contact, error)
}
