package mapping

import (
	"github.com/bitcoin-sv/spv-wallet/api"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/contacts/contactsmodels"
	"github.com/bitcoin-sv/spv-wallet/lox"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/samber/lo"
)

// MapToContactContract maps a contact to a response
func MapToContactContract(c *contactsmodels.Contact) api.ModelsContact {
	return api.ModelsContact{
		Id:        c.ID,
		FullName:  c.FullName,
		Paymail:   c.Paymail,
		PubKey:    c.PubKey,
		Status:    mapContactStatus(c.Status),
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
		DeletedAt: c.DeletedAt,
	}
}

// ContactsPagedResponse maps a paged result of contacts to a response.
func ContactsPagedResponse(contacts *models.PagedResult[contactsmodels.Contact]) api.ModelsContactsSearchResult {
	return api.ModelsContactsSearchResult{
		Page: api.ModelsSearchPage{
			Size:          contacts.PageDescription.Size,
			Number:        contacts.PageDescription.Number,
			TotalElements: contacts.PageDescription.TotalElements,
			TotalPages:    contacts.PageDescription.TotalPages,
		},
		Content: lo.Map(contacts.Content, lox.MappingFn(ContactsResponse)),
	}
}

// ContactsResponse maps an operation to a response.
func ContactsResponse(operation *contactsmodels.Contact) api.ModelsContact {
	return api.ModelsContact{
		Id:        operation.ID,
		FullName:  operation.FullName,
		Paymail:   operation.Paymail,
		PubKey:    operation.PubKey,
		Status:    mapContactStatus(operation.Status),
		UpdatedAt: operation.UpdatedAt,
		DeletedAt: operation.DeletedAt,
		CreatedAt: operation.CreatedAt,
	}
}

func mapContactStatus(s string) api.ModelsContactStatus {
	switch s {
	case contactsmodels.ContactNotConfirmed:
		return api.Unconfirmed
	case contactsmodels.ContactAwaitAccept:
		return api.Awaiting
	case contactsmodels.ContactConfirmed:
		return api.Confirmed
	case contactsmodels.ContactRejected:
		return api.Rejected
	default:
		return "unknown"
	}
}
