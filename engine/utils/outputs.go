package utils

import (
	"crypto/rand"
	"math"
	"math/big"
)

// SplitOutputValues splits the satoshis value randomly into nrOfOutputs pieces
func SplitOutputValues(satoshis uint64, nrOfOutputs int) ([]uint64, error) {
	outputValues := make([]uint64, nrOfOutputs)
	outputUsed := uint64(0)
	for i := 0; i < nrOfOutputs; i++ {
		a, err := rand.Int(
			rand.Reader, big.NewInt(math.MaxInt64),
		)
		if err != nil {
			return nil, err
		}
		randomOutput := (((float64(a.Int64()) / (1 << 63)) * 50) + 75) / 100
		outputValue := uint64(randomOutput * float64(satoshis) / float64(nrOfOutputs))

		outputValues[i] = outputValue
		outputUsed += outputValue
	}

	outputValues[len(outputValues)-1] += satoshis - outputUsed

	return outputValues, nil
}
