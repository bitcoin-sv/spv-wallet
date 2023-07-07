package mappings

import (
	"github.com/BuxOrg/bux"
	buxmodels "github.com/BuxOrg/bux-models"
)

func MapToUtxoPointer(u *bux.UtxoPointer) *buxmodels.UtxoPointer {
	return &buxmodels.UtxoPointer{
		TransactionID: u.TransactionID,
		OutputIndex:   (int)u.OutputIndex,
	}
}
