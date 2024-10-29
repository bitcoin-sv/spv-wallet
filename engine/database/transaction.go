package database

// Transaction represents a transaction in the database.
// Fixme: This is not integrated with out db engine yet.
type Transaction struct {
	ID       string `gorm:"type:char(64);primaryKey"`
	TxStatus TxStatus
}
