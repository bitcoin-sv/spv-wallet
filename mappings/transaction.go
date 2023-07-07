package mappings

import (
	"github.com/BuxOrg/bux"
	buxmodels "github.com/BuxOrg/bux-models"
	"github.com/BuxOrg/bux-server/mappings/common"
)

func MapToTransactionContract(t *bux.Transaction) *buxmodels.Transaction {
	return &buxmodels.Transaction{
		Model:                *common.MapToContract(&t.Model),
		ID:                   t.ID,
		Hex:                  t.Hex,
		XpubInIDs:            t.XpubInIDs,
		XpubOutIDs:           t.XpubOutIDs,
		BlockHash:            t.BlockHash,
		BlockHeight:          t.BlockHeight,
		Fee:                  t.Fee,
		NumberOfInputs:       t.NumberOfInputs,
		NumberOfOutputs:      t.NumberOfOutputs,
		DraftID:              t.DraftID,
		TotalValue:           t.TotalValue,
		OutputValue:          t.OutputValue,
		Status:               string(t.Status),
		TransactionDirection: string(t.Direction),
	}
}

func MapToTransactionConfigBux(tx *buxmodels.TransactionConfig) *bux.TransactionConfig {
	destinations := make([]*bux.Destination, 0)
	for _, destination := range tx.ChangeDestinations {
		destinations = append(destinations, MapToDestinationBux(destination))
	}

	return &bux.TransactionConfig{
		ChangeDestinations: destinations,
	}
}
