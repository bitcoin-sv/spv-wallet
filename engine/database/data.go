package database

// Data holds the data stored in outputs.
// Fixme: This is not integrated with out db engine yet.
type Data struct {
	TxID string `gorm:"primaryKey"`
	Vout uint32 `gorm:"primaryKey"`
	Blob []byte
}
