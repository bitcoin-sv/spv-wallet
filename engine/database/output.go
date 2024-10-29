package database

// Output represents an output of a transaction.
// Fixme: This is not integrated with out db engine yet.
type Output struct {
	TxID       string  `gorm:"primaryKey"`
	Vout       uint32  `gorm:"primaryKey"`
	SpendingTX *string `gorm:"type:char(64)"`
}

// IsSpent returns true if the output is spent.
func (o *Output) IsSpent() bool {
	return o.SpendingTX != nil
}

// Spend marks the output as spent.
func (o *Output) Spend(spendingTXID string) {
	o.SpendingTX = &spendingTXID
}
