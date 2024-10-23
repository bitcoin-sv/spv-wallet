package database

type Output struct {
	TxID       string  `gorm:"primaryKey"`
	Vout       uint32  `gorm:"primaryKey"`
	SpendingTX *string `gorm:"type:char(64)"`
}

func (o *Output) IsSpent() bool {
	return o.SpendingTX != nil
}

func (o *Output) Spend(spendingTXID string) {
	o.SpendingTX = &spendingTXID
}
