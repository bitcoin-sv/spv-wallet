package mappings

import (
	"github.com/BuxOrg/bux"
	spvwalletmodels "github.com/BuxOrg/bux-models"
	"github.com/BuxOrg/spv-wallet/mappings/common"
	customtypes "github.com/mrz1836/go-datastore/custom_types"
)

// MapToUtxoPointer will map the utxo-pointer model from spv-wallet to the spv-wallet-models contract
func MapToUtxoPointer(u *bux.UtxoPointer) *spvwalletmodels.UtxoPointer {
	if u == nil {
		return nil
	}

	return &spvwalletmodels.UtxoPointer{
		TransactionID: u.TransactionID,
		OutputIndex:   u.OutputIndex,
	}
}

// MapToUtxoPointerSPV will map the utxo-pointer model from spv-wallet-models to the spv-wallet contract
func MapToUtxoPointerSPV(u *spvwalletmodels.UtxoPointer) *bux.UtxoPointer {
	if u == nil {
		return nil
	}

	return &bux.UtxoPointer{
		TransactionID: u.TransactionID,
		OutputIndex:   u.OutputIndex,
	}
}

// MapToUtxoContract will map the utxo model from spv-wallet to the spv-wallet-models contract
func MapToUtxoContract(u *bux.Utxo) *spvwalletmodels.Utxo {
	if u == nil {
		return nil
	}

	return &spvwalletmodels.Utxo{
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

// MapToUtxoSPV will map the utxo model from spv-wallet-models to the spv-wallet contract
func MapToUtxoSPV(u *spvwalletmodels.Utxo) *bux.Utxo {
	if u == nil {
		return nil
	}

	var draftID customtypes.NullString
	draftID.String = u.DraftID

	var spendingTxID customtypes.NullString
	spendingTxID.String = u.SpendingTxID

	return &bux.Utxo{
		Model:        *common.MapToModel(&u.Model),
		UtxoPointer:  *MapToUtxoPointerSPV(&u.UtxoPointer),
		ID:           u.ID,
		XpubID:       u.XpubID,
		Satoshis:     u.Satoshis,
		ScriptPubKey: u.ScriptPubKey,
		Type:         u.Type,
		DraftID:      draftID,
		SpendingTxID: spendingTxID,
		Transaction:  MapToTransactionSPV(u.Transaction),
	}
}
