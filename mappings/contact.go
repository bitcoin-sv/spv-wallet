package mappings

import (
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/mappings/common"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

// MapToContactContract will map the contact to the spv-wallet-models contract
func MapToContactContract(src *engine.Contact) *response.Contact {
	if src == nil {
		return nil
	}

	return &response.Contact{
		ID:       src.ID,
		Model:    *common.MapToContract(&src.Model),
		FullName: src.FullName,
		Paymail:  src.Paymail,
		PubKey:   src.PubKey,
		Status:   mapContactStatus(src.Status),
	}
}

// MapToContactContracts will map the contacts collection to the spv-wallet-models contracts collection
func MapToContactContracts(src []*engine.Contact) []*response.Contact {
	res := make([]*response.Contact, 0, len(src))

	for _, c := range src {
		res = append(res, MapToContactContract(c))
	}

	return res
}

func mapContactStatus(s engine.ContactStatus) response.ContactStatus {
	switch s {
	case engine.ContactNotConfirmed:
		return response.ContactNotConfirmed
	case engine.ContactAwaitAccept:
		return response.ContactAwaitAccept
	case engine.ContactConfirmed:
		return response.ContactConfirmed
	case engine.ContactRejected:
		return response.ContactRejected
	default:
		return "unknown"
	}
}
