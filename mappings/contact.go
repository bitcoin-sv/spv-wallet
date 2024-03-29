package mappings

import (
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/mappings/common"
	"github.com/bitcoin-sv/spv-wallet/models"
)

// MapToContactContract will map the contact to the spv-wallet-models contract
func MapToContactContract(src *engine.Contact) *models.Contact {
	if src == nil {
		return nil
	}

	return &models.Contact{
		ID:       src.ID,
		Model:    *common.MapToContract(&src.Model),
		FullName: src.FullName,
		Paymail:  src.Paymail,
		PubKey:   src.PubKey,
		Status:   mapContactStatus(src.Status),
	}
}

func MapToContactContracts(src []*engine.Contact) []*models.Contact {
	res := make([]*models.Contact, len(src))

	for _, c := range src {
		res = append(res, MapToContactContract(c))
	}

	return res
}

func mapContactStatus(s engine.ContactStatus) string {
	switch s {
	case engine.ContactNotConfirmed:
		return "unconfirmed"
	case engine.ContactAwaitAccept:
		return "awaiting"
	case engine.ContactConfirmed:
		return "confirmed"
	case engine.ContactRejected:
		return "rejected"
	default:
		return "unknown"
	}
}
