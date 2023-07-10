package mappings

import (
	"github.com/BuxOrg/bux"
	buxmodels "github.com/BuxOrg/bux-models"
	"github.com/BuxOrg/bux-server/mappings/common"
)

// MapToXpubContract will map the xpub model from bux to the bux-models contract
func MapToXpubContract(xpub *bux.Xpub) *buxmodels.Xpub {
	return &buxmodels.Xpub{
		Model:           *common.MapToContract(&xpub.Model),
		ID:              xpub.ID,
		CurrentBalance:  xpub.CurrentBalance,
		NextInternalNum: xpub.NextInternalNum,
		NextExternalNum: xpub.NextExternalNum,
	}
}
