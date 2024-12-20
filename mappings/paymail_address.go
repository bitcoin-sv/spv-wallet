package mappings

import (
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/mappings/common"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

// MapToPaymailContract will map the spv-wallet paymail-address model to the spv-wallet-models contract
func MapToPaymailContract(pa *engine.PaymailAddress) *response.PaymailAddress {
	if pa == nil {
		return nil
	}

	return &response.PaymailAddress{
		Model:      *common.MapToContract(&pa.Model),
		ID:         pa.ID,
		XpubID:     pa.XpubID,
		Alias:      pa.Alias,
		Domain:     pa.Domain,
		PublicName: pa.PublicName,
		Avatar:     pa.Avatar,
		Address:    pa.String(),
	}
}

// MapToPaymailP4Contract will map the spv-wallet-models paymail-address contract to the spv-wallet paymail-address model
func MapToPaymailP4Contract(p *engine.PaymailP4) *response.PaymailP4 {
	if p == nil {
		return nil
	}

	return &response.PaymailP4{
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

// MapPaymailP4ModelToEngine will map the spv-wallet-models paymail-address contract to the spv-wallet paymail-address model
func MapPaymailP4ModelToEngine(p *response.PaymailP4) *engine.PaymailP4 {
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
