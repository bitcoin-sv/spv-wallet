package mappings

import (
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/models"
)

// MapToOldPaymailP4Contract will map the spv-wallet-models paymail-address contract to the spv-wallet paymail-address model
func MapToOldPaymailP4Contract(p *engine.PaymailP4) *models.PaymailP4 {
	if p == nil {
		return nil
	}

	return &models.PaymailP4{
		Alias:           p.Alias,
		Domain:          p.Domain,
		FromPaymail:     p.FromPaymail,
		Note:            p.Note,
		PubKey:          p.PubKey,
		ReceiveEndpoint: p.ReceiveEndpoint,
		ReferenceID:     p.ReferenceID,
		ResolutionType:  p.ResolutionType,
	}
}

// MapOldPaymailP4ModelToEngine will map the spv-wallet-models paymail-address contract to the spv-wallet paymail-address model
func MapOldPaymailP4ModelToEngine(p *models.PaymailP4) *engine.PaymailP4 {
	if p == nil {
		return nil
	}

	return &engine.PaymailP4{
		Alias:           p.Alias,
		Domain:          p.Domain,
		FromPaymail:     p.FromPaymail,
		Note:            p.Note,
		PubKey:          p.PubKey,
		ReceiveEndpoint: p.ReceiveEndpoint,
		ReferenceID:     p.ReferenceID,
		ResolutionType:  p.ResolutionType,
	}
}
