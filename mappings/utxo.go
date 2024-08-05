package mappings

import (
	"github.com/bitcoin-sv/spv-wallet/engine"
	customtypes "github.com/bitcoin-sv/spv-wallet/engine/datastore/customtypes"
	"github.com/bitcoin-sv/spv-wallet/mappings/common"
	"github.com/bitcoin-sv/spv-wallet/models"
)

// MapToUtxoPointer will map the utxo-pointer model from spv-wallet to the spv-wallet-models contract
func MapToUtxoPointer(u *engine.UtxoPointer) *models.UtxoPointer {
	if u == nil {
		return nil
	}

	return &models.UtxoPointer{
		TransactionID: u.TransactionID,
		OutputIndex:   u.OutputIndex,
	}
}

// MapUtxoPointerModelToEngine will map the utxo-pointer model from spv-wallet-models to the spv-wallet contract
func MapUtxoPointerModelToEngine(u *models.UtxoPointer) *engine.UtxoPointer {
	if u == nil {
		return nil
	}

	return &engine.UtxoPointer{
		TransactionID: u.TransactionID,
		OutputIndex:   u.OutputIndex,
	}
}

// MapToUtxoContract will map the utxo model from spv-wallet to the spv-wallet-models contract
func MapToUtxoContract(u *engine.Utxo) *models.Utxo {
	if u == nil {
		return nil
	}

	return &models.Utxo{
		Model:        *common.MapToOldContract(&u.Model),
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

// MapUtxoModelToEngine will map the utxo model from spv-wallet-models to the spv-wallet contract
func MapUtxoModelToEngine(u *models.Utxo) *engine.Utxo {
	if u == nil {
		return nil
	}

	var draftID customtypes.NullString
	draftID.String = u.DraftID

	var spendingTxID customtypes.NullString
	spendingTxID.String = u.SpendingTxID

	return &engine.Utxo{
		Model:        *common.MapOldContractToModel(&u.Model),
		UtxoPointer:  *MapUtxoPointerModelToEngine(&u.UtxoPointer),
		ID:           u.ID,
		XpubID:       u.XpubID,
		Satoshis:     u.Satoshis,
		ScriptPubKey: u.ScriptPubKey,
		Type:         u.Type,
		DraftID:      draftID,
		SpendingTxID: spendingTxID,
		Transaction:  MapTransactionModelToEngine(u.Transaction),
	}
}
