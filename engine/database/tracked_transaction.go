package database

import (
	trx "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// TrackedTransaction represents a transaction in the database.
type TrackedTransaction struct {
	ID       string `gorm:"type:char(64);primaryKey"`
	TxStatus TxStatus

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

func (t *TrackedTransaction) BeforeSave(tx *gorm.DB) (err error) {
	for _, output := range t.Outputs {
		if output.IsSpent() {
			return spverrors.Newf("output %s is already spent", output.Outpoint())
		}
		t.TrackedOutputs = append(t.TrackedOutputs, output.ToTrackedOutput())
	}
	return nil
}
