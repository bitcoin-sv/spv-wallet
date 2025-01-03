package database

import (
	trx "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"gorm.io/datatypes"
	"gorm.io/gorm"
	"slices"
	"time"
)

// TrackedTransaction represents a transaction in the database.
type TrackedTransaction struct {
	ID       string `gorm:"type:char(64);primaryKey"`
	TxStatus TxStatus

	CreatedAt time.Time
	UpdatedAt time.Time

	BUMP *datatypes.JSONType[trx.MerklePath]

	Data []*Data `gorm:"foreignKey:TxID"`

	Outputs       []Output         `gorm:"-"`
	TrackedInputs []*TrackedOutput `gorm:"foreignKey:SpendingTX"`

	// TrackedOutputs are automatically populated from Outputs and Inputs.
	TrackedOutputs []*TrackedOutput `gorm:"foreignKey:TxID"`
}

// AddOutputs adds outputs to the transaction.
func (t *TrackedTransaction) AddOutputs(outputs ...Output) {
	t.Outputs = append(t.Outputs, outputs...)
}

// AddInputs adds inputs to the transaction.
func (t *TrackedTransaction) AddInputs(inputs ...*TrackedOutput) {
	t.TrackedInputs = append(t.TrackedInputs, inputs...)
}

// AddData adds data to the transaction.
func (t *TrackedTransaction) AddData(data ...*Data) {
	t.Data = append(t.Data, data...)
}

func (t *TrackedTransaction) BeforeSave(tx *gorm.DB) error {
	for _, output := range t.Outputs {
		if output.IsSpent() {
			return spverrors.Newf("output %s is already spent", output.Outpoint())
		}
		t.TrackedOutputs = append(t.TrackedOutputs, output.ToTrackedOutput())
	}
	return nil
}

func (t *TrackedTransaction) AfterSave(tx *gorm.DB) error {
	// Add new UTXOs
	userUTXOs := slices.Collect(func(yield func(utxos *UserUtxos) bool) {
		for _, output := range t.Outputs {
			yield(output.ToUserUTXO())
		}
	})
	if len(userUTXOs) > 0 {
		err := tx.Model(&UserUtxos{}).Create(userUTXOs).Error
		if err != nil {
			return spverrors.Wrapf(err, "failed to save user utxos")
		}
	}

	// Remove spent UTXOs
	spentOutpoints := slices.Collect(func(yield func(sqlPair []any) bool) {
		for _, outpoint := range t.TrackedInputs {
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
