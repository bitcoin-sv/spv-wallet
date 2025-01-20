package inputs

import (
	"context"
	"time"

	"github.com/bitcoin-sv/spv-wallet/engine/database"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	txerrors "github.com/bitcoin-sv/spv-wallet/engine/transaction/errors"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"gorm.io/gorm"
)

const txIdColumn = "tx_id"
const voutColumn = "vout"

// Selector is a service that selects inputs for transaction.
type Selector interface {
	SelectInputsForTransaction(ctx context.Context, userID string, satoshis bsv.Satoshis, byteSizeOfTxBeforeAddingSelectedInputs uint64) ([]*database.UserUTXO, error)
}

const (
	// estimatedChangeOutputSize is the estimated size of a change output
	// that will be added to transaction in case there are a change from transaction.
	// Currently, for this estimation we're assuming single change output with P2PKH locking script.
	estimatedChangeOutputSize = 34
)

type sqlInputsSelector struct {
	feeUnit bsv.FeeUnit
	db      *gorm.DB
}

// NewSelector creates a new instance of Selector.
func NewSelector(db *gorm.DB, feeUnit bsv.FeeUnit) Selector {
	return &sqlInputsSelector{
		db:      db,
		feeUnit: feeUnit,
	}
}

func (r *sqlInputsSelector) SelectInputsForTransaction(ctx context.Context, userID string, outputsTotalValue bsv.Satoshis, byteSizeOfTxWithoutInputs uint64) (utxos []*database.UserUTXO, err error) {
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

func (r *sqlInputsSelector) buildQueryForInputs(db *gorm.DB, userID string, outputsTotalValue bsv.Satoshis, txWithoutInputsSize uint64) *gorm.DB {
	composer := &inputsQueryComposer{
		userID:              userID,
		outputsTotalValue:   outputsTotalValue,
		txWithoutInputsSize: txWithoutInputsSize,
		feeUnit:             r.feeUnit,
	}
	return composer.build(db)
}

func (r *sqlInputsSelector) buildUpdateTouchedAtQuery(db *gorm.DB, utxos []*database.UserUTXO) *gorm.DB {
	outpoints := make([][]any, 0, len(utxos))
	for _, utxo := range utxos {
		outpoints = append(outpoints, []any{utxo.TxID, utxo.Vout})
	}
	return db.Model(&database.UserUTXO{}).Where("(tx_id, vout) in (?)", outpoints)
}
