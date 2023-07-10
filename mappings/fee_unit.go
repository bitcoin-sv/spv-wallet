package mappings

import (
	buxmodels "github.com/BuxOrg/bux-models"
	"github.com/BuxOrg/bux/utils"
)

// MapToFeeUnitContract will map the fee-unit model from bux to the bux-models contract
func MapToFeeUnitContract(fu *utils.FeeUnit) (fc *buxmodels.FeeUnit) {
	return &buxmodels.FeeUnit{
		Satoshis: fu.Satoshis,
		Bytes:    fu.Bytes,
	}
}

// MapToFeeUnitBux will map the fee-unit model from bux-models to the bux contract
func MapToFeeUnitBux(fu *buxmodels.FeeUnit) (fc *utils.FeeUnit) {
	return &utils.FeeUnit{
		Satoshis: fu.Satoshis,
		Bytes:    fu.Bytes,
	}
}
