package mappings

import (
	"github.com/bitcoin-sv/spv-wallet/engine/utils"
	"github.com/bitcoin-sv/spv-wallet/models"
)

// MapToFeeUnitContract will map the fee-unit model from spv-wallet to the spv-wallet-models contract
func MapToFeeUnitContract(fu *utils.FeeUnit) (fc *models.FeeUnit) {
	if fu == nil {
		return nil
	}

	return &models.FeeUnit{
		Satoshis: fu.Satoshis,
		Bytes:    fu.Bytes,
	}
}

// MapToFeeUnitSPV will map the fee-unit model from spv-wallet-models to the spv-wallet contract
func MapToFeeUnitSPV(fu *models.FeeUnit) (fc *utils.FeeUnit) {
	if fu == nil {
		return nil
	}

	return &utils.FeeUnit{
		Satoshis: fu.Satoshis,
		Bytes:    fu.Bytes,
	}
}
