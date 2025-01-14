package database

import (
	"slices"
	"time"

	trx "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// TrackedTransaction represents a transaction in the database.
type TrackedTransaction struct {
	ID       string `gorm:"type:char(64);primaryKey"`
	TxStatus TxStatus

	CreatedAt time.Time
	UpdatedAt time.Time

	BUMP *datatypes.JSONType[trx.MerklePath]

	Data []*Data `gorm:"foreignKey:TxID"`

	Inputs  []*Output `gorm:"foreignKey:SpendingTX"`
	Outputs []*Output `gorm:"foreignKey:TxID"`

	newUTXOs []*UserUtxos `gorm:"-"`
}

func (t *TrackedTransaction) CreateP2PKHOutput(output *Output, customInstructions datatypes.JSONSlice[CustomInstruction]) {
	t.Outputs = append(t.Outputs, output)
	t.newUTXOs = append(t.newUTXOs, NewP2PKHUserUTXO(output, customInstructions))
}

func (t *TrackedTransaction) CreateDataOutput(data *Data, userID string) {
	t.Data = append(t.Data, data) //TODO: Most probably Data should be also linked to the user
	t.Outputs = append(t.Outputs, &Output{
		TxID:     data.TxID,
		Vout:     data.Vout,
		UserID:   userID,
		Satoshis: 0,
	})
}

// AddInputs adds inputs to the transaction.
func (t *TrackedTransaction) AddInputs(inputs ...*Output) {
	t.Inputs = append(t.Inputs, inputs...)
}

// AfterCreate is a hook that is called after creating the transaction.
// It is responsible for adding new (User's) UTXOs and removing spent UTXOs.
func (t *TrackedTransaction) AfterCreate(tx *gorm.DB) error {
	// Add new UTXOs
	if len(t.newUTXOs) > 0 {
		err := tx.Model(&UserUtxos{}).Create(t.newUTXOs).Error
		if err != nil {
			return spverrors.Wrapf(err, "failed to save user utxos")
		}
	}

	// Remove spent UTXOs
	spentOutpoints := slices.AppendSeq(
		make([][]any, 0, len(t.Inputs)),
		func(yield func(sqlPair []any) bool) {
			for _, outpoint := range t.Inputs {
				yield([]any{outpoint.TxID, outpoint.Vout})
			}
		})
	if len(spentOutpoints) > 0 {
		err := tx.Where("(tx_id, vout) IN ?", spentOutpoints).Delete(&UserUtxos{}).Error
		if err != nil {
			return spverrors.Wrapf(err, "failed to delete spent utxos")
		}
	}

	return nil
}
