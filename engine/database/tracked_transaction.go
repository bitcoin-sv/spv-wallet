package database

import (
	trx "github.com/bitcoin-sv/go-sdk/transaction"
	"gorm.io/datatypes"
)

// TrackedTransaction represents a transaction in the database.
type TrackedTransaction struct {
	ID       string `gorm:"type:char(64);primaryKey"`
	TxStatus TxStatus

	BUMP *datatypes.JSONType[trx.MerklePath]

	Outputs []*TrackedOutput `gorm:"foreignKey:TxID"`
	Data    []*Data          `gorm:"foreignKey:TxID"`
	Inputs  []*TrackedOutput `gorm:"foreignKey:SpendingTX"`
}

// AddOutputs adds outputs to the transaction.
func (t *TrackedTransaction) AddOutputs(outputs ...*TrackedOutput) {
	t.Outputs = append(t.Outputs, outputs...)
}

// AddInputs adds inputs to the transaction.
func (t *TrackedTransaction) AddInputs(inputs ...*TrackedOutput) {
	t.Inputs = append(t.Inputs, inputs...)
}

// AddData adds data to the transaction.
func (t *TrackedTransaction) AddData(data ...*Data) {
	t.Data = append(t.Data, data...)
}
