package outlines_test

import (
	"context"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/outlines/testabilities"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/bitcoin-sv/spv-wallet/models/transaction/bucket"
)

const (
	//this datasize will result in size 999 and fee 1 (before adding change output) and size 1033 and fee 2 (after adding change output)
	higherDatasizeForMinimumFee = 825
	minimumDatasizeForFee2      = higherDatasizeForMinimumFee + 1
)

func TestOutlineWithChange(t *testing.T) {
	tests := map[string]struct {
		utxoValue      bsv.Satoshis
		datasize       uint64
		expectedChange bsv.Satoshis
	}{
		"standard case": {
			utxoValue:      10,
			datasize:       100,
			expectedChange: 9,
		},
		"standard case for higher fee": {
			utxoValue:      10,
			datasize:       minimumDatasizeForFee2,
			expectedChange: 8,
		},
		"addition of change output makes higher fee": {
			utxoValue:      10,
			datasize:       higherDatasizeForMinimumFee,
			expectedChange: 8,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			given, then := testabilities.New(t)

			// given:
			service := given.NewTransactionOutlinesService()

			// and:
			given.UTXOSelector().WillReturnUTXOs(test.utxoValue)

			// when:
			tx, err := service.CreateBEEF(context.Background(), given.TransactionSpecWithDatasize(test.datasize))

			// then:
			thenTx := then.Created(tx).WithNoError(err).WithParseableBEEFHex()

			thenTx.HasOutputs(2)

			thenTx.Output(1).
				HasBucket(bucket.BSV).
				HasSatoshis(test.expectedChange).
				UnlockableBySender()
		})
	}
}

func TestOutlineNoChange(t *testing.T) {
	tests := map[string]struct {
		utxoValue bsv.Satoshis
		datasize  uint64
	}{
		"standard case": {
			utxoValue: 1,
			datasize:  100,
		},
		"standard case for higher fee": {
			utxoValue: 2,
			datasize:  minimumDatasizeForFee2,
		},
		"addition of change output makes higher fee and change decreases to zero": {
			datasize:  higherDatasizeForMinimumFee,
			utxoValue: 2,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			given, then := testabilities.New(t)

			// given:
			service := given.NewTransactionOutlinesService()

			// and:
			given.UTXOSelector().WillReturnUTXOs(test.utxoValue)

			// when:
			tx, err := service.CreateBEEF(context.Background(), given.TransactionSpecWithDatasize(test.datasize))

			// then:
			thenTx := then.Created(tx).WithNoError(err).WithParseableBEEFHex()

			thenTx.HasOutputs(1)
		})
	}
}
