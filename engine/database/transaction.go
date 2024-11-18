package database

// Transaction represents a transaction in the database.
type Transaction struct {
	ID       string `gorm:"type:char(64);primaryKey"`
	TxStatus TxStatus

	Outputs []*Output `gorm:"foreignKey:TxID"`
	Data    []*Data   `gorm:"foreignKey:TxID"`
}

// TableName implements gorm.Tabler to override automatic table naming.
// NOTE: This is because we have already a legacy table named "transactions".
// TODO: Remove this when we have migrated all data.
func (t *Transaction) TableName() string {
	return "new_transactions"
}

// AddOutputs adds outputs to the transaction.
func (t *Transaction) AddOutputs(outputs ...*Output) {
	t.Outputs = append(t.Outputs, outputs...)
}

// AddData adds data to the transaction.
func (t *Transaction) AddData(data ...*Data) {
	t.Data = append(t.Data, data...)
}
