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
