package mappings

import (
	buxmodels "github.com/BuxOrg/bux-models"
	"github.com/BuxOrg/bux/utils"
)

func MapToFeeUnitContract(fu *utils.FeeUnit) (fc *buxmodels.FeeUnit) {
	return &buxmodels.FeeUnit{
		Satoshis: fu.Satoshis,
		Bytes:    fu.Bytes,
	}
}

func MapToFeeUnitBux(fu *buxmodels.FeeUnit) (fc *utils.FeeUnit) {
	return &utils.FeeUnit{
		Satoshis: fu.Satoshis,
		Bytes:    fu.Bytes,
	}
}
