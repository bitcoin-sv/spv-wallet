package bsv

import (
	"strings"

	"github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/bsv/bsverrors"
)

type TxHexFormat string

const (
	TxHexFormatBEEF TxHexFormat = "BEEF"
	TxHexFormatRAW  TxHexFormat = "RAW"
)

func ParseTxHexFormat(s string) (TxHexFormat, error) {
	switch strings.ToUpper(s) {
	case "BEEF":
		return TxHexFormatBEEF, nil
	case "RAW":
		return TxHexFormatRAW, nil
	default:
		return "", bsverrors.ErrUnknownTransactionFormat
	}
}

// TxHex is a hex representation of a transaction.
type TxHex string

// IsBEEF checks if the transaction hex is a BEEF hex.
func (h TxHex) IsBEEF() bool {
	return strings.HasPrefix(string(h), "0100BEEF") || strings.HasPrefix(string(h), "0100beef")
}

// IsRawTx checks if the transaction hex is a raw transaction hex.
func (h TxHex) IsRawTx() bool {
	return !h.IsBEEF()
}

// ToBEEFTransaction converts the transaction hex to a BEEF transaction.
func (h TxHex) ToBEEFTransaction() (*transaction.Transaction, error) {
	if !h.IsBEEF() {
		return nil, spverrors.Newf("transaction hex is not a BEEF hex")
	}
	return transaction.NewTransactionFromBEEFHex(string(h)) //nolint:wrapcheck // we will handle this error in upper layers
}

func (h TxHex) ToRawTransaction() (*transaction.Transaction, error) {
	if !h.IsRawTx() {
		return nil, spverrors.Newf("transaction hex is not a raw hex")
	}
	return transaction.NewTransactionFromHex(string(h)) //nolint:wrapcheck // we will handle this error in upper layers
}

// Format returns the name of the format of the transaction hex.
func (h TxHex) Format() TxHexFormat {
	if h.IsBEEF() {
		return TxHexFormatBEEF
	}
	return TxHexFormatRAW
}
