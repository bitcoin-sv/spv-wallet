package database

type Data struct {
	TxID string `gorm:"primaryKey"`
	Vout uint32 `gorm:"primaryKey"`
	Blob []byte
}
