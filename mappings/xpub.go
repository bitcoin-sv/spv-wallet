package mappings

import (
	"github.com/bitcoin-sv/spv-wallet/engine"
	"github.com/bitcoin-sv/spv-wallet/mappings/common"
	spvwalletmodels "github.com/bitcoin-sv/spv-wallet/models"
)

// MapToXpubContract will map the xpub model from spv-wallet to the spv-wallet-models contract
func MapToXpubContract(xpub *engine.Xpub) *spvwalletmodels.Xpub {
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
