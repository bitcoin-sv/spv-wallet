package sql

import (
	"context"
	"time"

	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/database"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/errors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/outlines"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"gorm.io/gorm"
)

const (
	txIdColumn               = "tx_id"
	voutColumn               = "vout"
	minChange                = "min_change"
	customInstructionsColumn = "custom_instructions"
)

const (
	// estimatedChangeOutputSize is the estimated size of a change output
	// that will be added to transaction in case there are a change from transaction.
	// Currently, for this estimation we're assuming single change output with P2PKH locking script.
	estimatedChangeOutputSize = 34
)

// UTXOSelector is responsible for selecting UTXOs for a transaction in SQL databases.
type UTXOSelector struct {
	feeUnit bsv.FeeUnit
	db      *gorm.DB
}

// NewUTXOSelector creates a new instance of UTXOSelector.
func NewUTXOSelector(db *gorm.DB, feeUnit bsv.FeeUnit) *UTXOSelector {
	return &UTXOSelector{
		db:      db,
		feeUnit: feeUnit,
	}
}

// Select selects UTXOs of user to fund a transaction.
func (r *UTXOSelector) Select(ctx context.Context, tx *sdk.Transaction, userID string) (utxos []*outlines.UTXO, change bsv.Satoshis, err error) {
	// NOTE: this approach assumes that tx doesn't contain any predefined inputs and all should be selected to cover outputs
	outputsTotalValue := tx.TotalOutputSatoshis()
	byteSizeOfTxToFund := outputOnlyTxSize(tx.Outputs)

	var selected []*selectedUTXO
	selected, err = r.selectInputsForTransaction(ctx, userID, bsv.Satoshis(outputsTotalValue), byteSizeOfTxToFund)
	if err != nil {
		return nil, bsv.Satoshis(0), err
	}

	if len(selected) > 0 {
		// final change value, calculated by SQL, is present in all rows
		change = bsv.Satoshis(selected[0].Change)
	}

	utxos = make([]*outlines.UTXO, len(selected))
	for i, utxo := range selected {
		utxos[i] = &outlines.UTXO{
			TxID:               utxo.TxID,
			Vout:               utxo.Vout,
			CustomInstructions: bsv.CustomInstructions(utxo.CustomInstructions),
		}
	}
	return
}

func (r *UTXOSelector) selectInputsForTransaction(ctx context.Context, userID string, outputsTotalValue bsv.Satoshis, byteSizeOfTxWithoutInputs uint64) (utxos []*selectedUTXO, err error) {
	err = r.db.WithContext(ctx).Transaction(func(db *gorm.DB) error {
		inputsQuery := r.buildQueryForInputs(db, userID, outputsTotalValue, byteSizeOfTxWithoutInputs)

		if err := inputsQuery.Find(&utxos).Error; err != nil {
			utxos = nil
			return spverrors.Wrapf(err, "failed to select utxos for transaction")
		}

		if len(utxos) == 0 {
			return nil
		}

		updateQuery := r.buildUpdateTouchedAtQuery(db, utxos)

		if err := updateQuery.Update("touched_at", time.Now()).Error; err != nil {
			utxos = nil
			return spverrors.Wrapf(err, "failed to update touched_at for selected inputs")
		}

		return nil
	})
	if err != nil {
		return nil, txerrors.ErrUnexpectedErrorDuringInputsSelection.Wrap(err)
	}

	return utxos, nil
}

func (r *UTXOSelector) buildQueryForInputs(db *gorm.DB, userID string, outputsTotalValue bsv.Satoshis, txWithoutInputsSize uint64) *gorm.DB {
	composer := &inputsQueryComposer{
		userID:              userID,
		outputsTotalValue:   outputsTotalValue,
		txWithoutInputsSize: txWithoutInputsSize,
		feeUnit:             r.feeUnit,
	}
	return composer.build(db)
}

func (r *UTXOSelector) buildUpdateTouchedAtQuery(db *gorm.DB, utxos []*selectedUTXO) *gorm.DB {
	outpoints := make([][]any, 0, len(utxos))
	for _, utxo := range utxos {
		outpoints = append(outpoints, []any{utxo.TxID, utxo.Vout})
	}
	return db.Model(&database.UserUTXO{}).Where("(tx_id, vout) in (?)", outpoints)
}
