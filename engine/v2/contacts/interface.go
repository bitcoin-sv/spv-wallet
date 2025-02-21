package contacts

import (
	"context"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/contacts/contactsmodels"
)

// ContactRepo is a contacts repository
type ContactRepo interface {
	Create(ctx context.Context, newContact *contactsmodels.NewContact) error
	Update(ctx context.Context, contact *contactsmodels.NewContact) error
	Find(ctx context.Context, userID, paymail string) (*contactsmodels.Contact, error)
}
