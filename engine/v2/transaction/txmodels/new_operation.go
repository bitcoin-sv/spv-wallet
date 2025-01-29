package txmodels

import (
	"github.com/bitcoin-sv/spv-wallet/conv"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
)

// NewOperation is a data for creating a new operation.
type NewOperation struct {
	UserID string

	Counterparty string
	Type         string
	Value        int64

	Transaction *NewTransaction
}

// Add adds satoshis to the operation.
func (o *NewOperation) Add(satoshi bsv.Satoshis) {
	signedSatoshi, err := conv.Uint64ToInt64(uint64(satoshi))
	if err != nil {
		panic(err)
	}
	o.Value = o.Value + signedSatoshi
}

// Subtract subtracts satoshis from the operation.
func (o *NewOperation) Subtract(satoshi bsv.Satoshis) {
	signedSatoshi, err := conv.Uint64ToInt64(uint64(satoshi))
	if err != nil {
		panic(err)
	}
	o.Value = o.Value - signedSatoshi
}
