package mappings

import (
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/mappings/common"
	"github.com/bitcoin-sv/spv-wallet/models"
)

// MapToContactContract will map the contact to the spv-wallet-models contract
func MapToContactContract(c *engine.Contact) *models.Contact {
	if c == nil {
		return nil
	}

	return &models.Contact{
		ID:       c.ID,
		Model:    *common.MapToContract(&c.Model),
		FullName: c.FullName,
		Paymail:  c.Paymail,
		PubKey:   c.PubKey,
		Status:   mapContactStatus(c.Status),
	}
}

func mapContactStatus(s engine.ContactStatus) string {
	switch s {
	case engine.ContactNotConfirmed:
		return "unconfirmed"
	case engine.ContactAwaitAccept:
		return "awaiting"
	case engine.ContactConfirmed:
		return "confirmed"
	default:
		return "unknown"
	}
}
