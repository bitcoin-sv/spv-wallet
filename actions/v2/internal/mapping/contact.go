package mapping

import (
	"github.com/bitcoin-sv/spv-wallet/api"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/contacts/contactsmodels"
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
