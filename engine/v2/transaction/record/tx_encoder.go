package record

import (
	trx "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
)

type txEncoder struct {
	tx *trx.Transaction
}

func newTxEncoder(tx *trx.Transaction) *txEncoder {
	return &txEncoder{tx: tx}
}

func (e *txEncoder) ToBEEF() (string, error) {
	hex, err := e.tx.BEEFHex()
	if err != nil {
		return "", spverrors.Wrapf(err, "failed to encode transaction to BEEF")
	}
	return hex, nil
}

func (e *txEncoder) ToRawHEX() string {
	return e.tx.Hex()
}
