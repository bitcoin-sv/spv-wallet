package mappings

import (
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/mappings/common"
	spvwalletmodels "github.com/bitcoin-sv/spv-wallet/models"
)

// MapToPaymailContract will map the spv-wallet paymail-address model to the spv-wallet-models contract
func MapToPaymailContract(pa *engine.PaymailAddress) *spvwalletmodels.PaymailAddress {
	if pa == nil {
		return nil
	}

	return &spvwalletmodels.PaymailAddress{
		Model:           *common.MapToContract(&pa.Model),
		ID:              pa.ID,
		XpubID:          pa.XpubID,
		Alias:           pa.Alias,
		Domain:          pa.Domain,
		PublicName:      pa.PublicName,
		Avatar:          pa.Avatar,
		ExternalXpubKey: pa.ExternalXpubKey,
	}
}

// MapToPaymailP4Contract will map the spv-wallet-models paymail-address contract to the spv-wallet paymail-address model
func MapToPaymailP4Contract(p *engine.PaymailP4) *spvwalletmodels.PaymailP4 {
	if p == nil {
		return nil
	}

	return &spvwalletmodels.PaymailP4{
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

// MapToPaymailP4SPV will map the spv-wallet-models paymail-address contract to the spv-wallet paymail-address model
func MapToPaymailP4SPV(p *spvwalletmodels.PaymailP4) *engine.PaymailP4 {
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
