package mappings

import (
	"github.com/BuxOrg/bux"
	spvwalletmodels "github.com/BuxOrg/bux-models"
	"github.com/BuxOrg/spv-wallet/mappings/common"
)

// MapToPaymailContract will map the spv-wallet paymail-address model to the spv-wallet-models contract
func MapToPaymailContract(pa *bux.PaymailAddress) *spvwalletmodels.PaymailAddress {
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
func MapToPaymailP4Contract(p *bux.PaymailP4) *spvwalletmodels.PaymailP4 {
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
func MapToPaymailP4SPV(p *spvwalletmodels.PaymailP4) *bux.PaymailP4 {
	if p == nil {
		return nil
	}

	return &bux.PaymailP4{
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
