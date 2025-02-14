package outlines

import (
	"math"

	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
)

const (
	txEnvelopeSize = 8 // version + locktime
)

func calculateFee(inputs annotatedInputs, outputs annotatedOutputs, feeUnit bsv.FeeUnit) bsv.Satoshis {
	size := estimatedSize(inputs, outputs)

	chunks := uint64(math.Ceil(float64(size) / float64(feeUnit.Bytes)))
	return bsv.Satoshis(chunks) * feeUnit.Satoshis
}

func estimatedSize(inputs annotatedInputs, outputs annotatedOutputs) uint64 {
	var size uint64

	size += txEnvelopeSize
	size += estimatedInputsSize(inputs)
	size += outputsSize(outputs)

	return size
}

func outputsSize(outputs annotatedOutputs) uint64 {
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

func estimatedInputsSize(inputs annotatedInputs) uint64 {
	var size uint64

	// input count:
	size += varIntSize(len(inputs))

	// inputs:
	for _, in := range inputs {
		size += in.estimatedSize
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
