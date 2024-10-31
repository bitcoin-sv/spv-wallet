package database

import "github.com/bitcoin-sv/spv-wallet/models/bsv"

// Data holds the data stored in outputs.
// Fixme: This is not integrated with out db engine yet.
type Data struct {
	TxID string `gorm:"primaryKey"`
	Vout uint32 `gorm:"primaryKey"`
	Blob []byte
}

// Outpoint returns bsv.Outpoint object which identifies the data-output.
func (o *Data) Outpoint() *bsv.Outpoint {
	return &bsv.Outpoint{
		TxID: o.TxID,
		Vout: o.Vout,
	}
}
