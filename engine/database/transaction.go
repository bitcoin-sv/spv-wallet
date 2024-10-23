package database

type Transaction struct {
	ID       string `gorm:"type:char(64);primaryKey"`
	TxStatus TxStatus
}
