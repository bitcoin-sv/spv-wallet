package datamodels

import "github.com/bitcoin-sv/spv-wallet/models/bsv"

// Data is a domain model for data stored in outputs (e.g. OP_RETURN).
type Data struct {
	TxID string
	Vout uint32

	UserID string

	Blob []byte
}

// ID returns the unique identifier of the data (outpoint string)
func (d *Data) ID() string {
	return bsv.Outpoint{TxID: d.TxID, Vout: d.Vout}.String()
}
