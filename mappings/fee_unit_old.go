package mappings

import (
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
)

// MapToOldFeeUnitContract will map the fee-unit model from spv-wallet to the spv-wallet-models contract
func MapToOldFeeUnitContract(fu *bsv.FeeUnit) (fc *models.FeeUnit) {
	if fu == nil {
		return nil
	}

	return &models.FeeUnit{
		Satoshis: int(fu.Satoshis), //nolint:gosec
		Bytes:    fu.Bytes,
	}
}

// MapOldFeeUnitModelToEngine will map the fee-unit model from spv-wallet-models to the spv-wallet contract
func MapOldFeeUnitModelToEngine(fu *models.FeeUnit) (fc *bsv.FeeUnit) {
	if fu == nil {
		return nil
	}

	return &bsv.FeeUnit{
		Satoshis: bsv.Satoshis(fu.Satoshis), //nolint:gosec
		Bytes:    fu.Bytes,
	}
}
