package database

import "github.com/bitcoin-sv/spv-wallet/models/bsv"

// Data holds the data stored in outputs.
type Data struct {
	TxID   string `gorm:"primaryKey"`
	Vout   uint32 `gorm:"primaryKey"`
	XpubID string
	Blob   []byte
}

// Outpoint returns bsv.Outpoint object which identifies the data-output.
func (o *Data) Outpoint() *bsv.Outpoint {
	return &bsv.Outpoint{
		TxID: o.TxID,
		Vout: o.Vout,
	}
}
