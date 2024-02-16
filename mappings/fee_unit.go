package mappings

import (
	spvwalletmodels "github.com/bitcoin-sv/bux-models"
	"github.com/bitcoin-sv/bux/utils"
)

// MapToFeeUnitContract will map the fee-unit model from spv-wallet to the spv-wallet-models contract
func MapToFeeUnitContract(fu *utils.FeeUnit) (fc *spvwalletmodels.FeeUnit) {
	if fu == nil {
		return nil
	}

	return &spvwalletmodels.FeeUnit{
		Satoshis: fu.Satoshis,
		Bytes:    fu.Bytes,
	}
}

// MapToFeeUnitSPV will map the fee-unit model from spv-wallet-models to the spv-wallet contract
func MapToFeeUnitSPV(fu *spvwalletmodels.FeeUnit) (fc *utils.FeeUnit) {
	if fu == nil {
		return nil
	}

	return &utils.FeeUnit{
		Satoshis: fu.Satoshis,
		Bytes:    fu.Bytes,
	}
}
