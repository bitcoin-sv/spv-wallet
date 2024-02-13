package mappings

import (
	"github.com/BuxOrg/bux"
	spvwalletmodels "github.com/BuxOrg/bux-models"
	"github.com/BuxOrg/spv-wallet/mappings/common"
)

// MapToXpubContract will map the xpub model from spv-wallet to the spv-wallet-models contract
func MapToXpubContract(xpub *bux.Xpub) *spvwalletmodels.Xpub {
	if xpub == nil {
		return nil
	}

	return &spvwalletmodels.Xpub{
		Model:           *common.MapToContract(&xpub.Model),
		ID:              xpub.ID,
		CurrentBalance:  xpub.CurrentBalance,
		NextInternalNum: xpub.NextInternalNum,
		NextExternalNum: xpub.NextExternalNum,
	}
}
