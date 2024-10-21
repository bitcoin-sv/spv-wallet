package mappings

import (
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/bitcoin-sv/spv-wallet/models/response"
)

// MapToFeeUnitContract will map the fee-unit model from spv-wallet to the spv-wallet-models contract
func MapToFeeUnitContract(fu *bsv.FeeUnit) (fc *response.FeeUnit) {
	if fu == nil {
		return nil
	}

	santoshis := int(fu.Satoshis) //nolint:gosec
	return &response.FeeUnit{
		Satoshis: santoshis,
		Bytes:    fu.Bytes,
	}
}

// MapFeeUnitModelToEngine will map the fee-unit model from spv-wallet-models to the spv-wallet contract
func MapFeeUnitModelToEngine(fu *response.FeeUnit) (fc *bsv.FeeUnit) {
	if fu == nil {
		return nil
	}

	return &bsv.FeeUnit{
		Satoshis: bsv.Satoshis(fu.Satoshis), //nolint:gosec
		Bytes:    fu.Bytes,
	}
}
