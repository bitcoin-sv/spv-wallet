package contacts

import (
	"context"

	"github.com/bitcoin-sv/spv-wallet/engine/v2/contacts/contactsmodels"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/filter"
)

// ContactRepo is a contacts repository
type ContactRepo interface {
	Create(ctx context.Context, newContact contactsmodels.NewContact) (*contactsmodels.Contact, error)
	Update(ctx context.Context, contact contactsmodels.NewContact) (*contactsmodels.Contact, error)
	UpdateStatus(ctx context.Context, userID, paymail, status string) error
	UpdateStatusByID(ctx context.Context, contactID uint, status string) (*contactsmodels.Contact, error)
	UpdateByID(ctx context.Context, contactID uint, fullName string) (*contactsmodels.Contact, error)
	Delete(ctx context.Context, userID, paymail string) error
	DeleteByID(ctx context.Context, contactID uint) error
	Find(ctx context.Context, userID, paymail string) (*contactsmodels.Contact, error)
	FindByID(ctx context.Context, contactID uint) (*contactsmodels.Contact, error)
	PaginatedForUser(ctx context.Context, userID string, page filter.Page, conditions map[string]interface{}) (*models.PagedResult[contactsmodels.Contact], error)
	PaginatedForAdmin(ctx context.Context, page filter.Page, conditions map[string]interface{}) (*models.PagedResult[contactsmodels.Contact], error)
}
