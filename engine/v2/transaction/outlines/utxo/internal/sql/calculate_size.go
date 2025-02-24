package sql

import (
	sdk "github.com/bitcoin-sv/go-sdk/transaction"
)

const (
	txEnvelopeSize = 8 // version + locktime
)

func outputOnlyTxSize(outputs []*sdk.TransactionOutput) uint64 {
	var size uint64
	size += txEnvelopeSize
	size += varIntSize(0) // inputs count
	size += outputsSize(outputs)
	return size
}

func outputsSize(outputs []*sdk.TransactionOutput) uint64 {
	var size uint64

	// output count:
	size += varIntSize(len(outputs))

	// outputs:
	for _, out := range outputs {
		size += 8
		scriptLen := len(*out.LockingScript)
		size += varIntSize(scriptLen) + toU64(scriptLen)
	}

	return size
}

//nolint:gosec // No need to check for overflows from int to uint64 here
func varIntSize(val int) uint64 {
	length := sdk.VarInt(val).Length()
	return toU64(length)
}

//nolint:gosec // No need to check for overflows from int to uint64 here
func toU64(val int) uint64 {
	return uint64(val)
}
