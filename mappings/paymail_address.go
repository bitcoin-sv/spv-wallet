package mappings

import (
	"github.com/BuxOrg/bux"
	buxmodels "github.com/BuxOrg/bux-models"
	"github.com/BuxOrg/bux-server/mappings/common"
)

func MapToPaymailContract(pa *bux.PaymailAddress) *buxmodels.PaymailAddress {
	return &buxmodels.PaymailAddress{
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

func MapToPaymailP4Contract(p *bux.PaymailP4) *buxmodels.PaymailP4 {
	return &buxmodels.PaymailP4{
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

func MapToPaymailP4Bux(p *buxmodels.PaymailP4) *bux.PaymailP4 {
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
