package mappings

import (
	"github.com/BuxOrg/bux"
	buxmodels "github.com/BuxOrg/bux-models"
	"github.com/BuxOrg/bux-server/mappings/common"
	customtypes "github.com/mrz1836/go-datastore/custom_types"
)

func MapToUtxoPointer(u *bux.UtxoPointer) *buxmodels.UtxoPointer {
	return &buxmodels.UtxoPointer{
		TransactionID: u.TransactionID,
		OutputIndex:   u.OutputIndex,
	}
}

func MapToUtxoPointerBux(u *buxmodels.UtxoPointer) *bux.UtxoPointer {
	return &bux.UtxoPointer{
		TransactionID: u.TransactionID,
		OutputIndex:   u.OutputIndex,
	}
}

func MapToUtxoContract(u *bux.Utxo) *buxmodels.Utxo {
	return &buxmodels.Utxo{
		Model:        *common.MapToContract(&u.Model),
		UtxoPointer:  *MapToUtxoPointer(&u.UtxoPointer),
		ID:           u.ID,
		XpubID:       u.XpubID,
		Satoshis:     u.Satoshis,
		ScriptPubKey: u.ScriptPubKey,
		Type:         u.Type,
		DraftID:      u.DraftID.String,
		SpendingTxID: u.SpendingTxID.String,
		Transaction:  MapToTransactionContract(u.Transaction),
	}
}

func MapToUtxoBux(u *buxmodels.Utxo) *bux.Utxo {
	var draftID customtypes.NullString
	draftID.String = u.DraftID

	var spendingTxID customtypes.NullString
	spendingTxID.String = u.SpendingTxID

	return &bux.Utxo{
		Model:        *common.MapToModel(&u.Model),
		UtxoPointer:  *MapToUtxoPointerBux(&u.UtxoPointer),
		ID:           u.ID,
		XpubID:       u.XpubID,
		Satoshis:     u.Satoshis,
		ScriptPubKey: u.ScriptPubKey,
		Type:         u.Type,
		DraftID:      draftID,
		SpendingTxID: spendingTxID,
		Transaction:  MapToTransactionBux(u.Transaction),
	}
}
