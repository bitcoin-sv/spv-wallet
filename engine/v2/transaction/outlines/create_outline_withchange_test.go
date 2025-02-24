package outlines_test

import (
	"context"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/outlines/testabilities"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/bitcoin-sv/spv-wallet/models/transaction/bucket"
)

func TestOutlineWithChange(t *testing.T) {
	given, then := testabilities.New(t)

	// given:
	service := given.NewTransactionOutlinesService()

	// and:
	change := bsv.Satoshis(9)
	utxoValue := bsv.Satoshis(10)

	// and:
	given.UTXOSelector().WillReturnUTXOs(change, utxoValue)

	// when:
	tx, err := service.CreateBEEF(context.Background(), given.MinimumValidTransactionSpec())

	// then:
	thenTx := then.Created(tx).WithNoError(err).WithParseableBEEFHex()

	thenTx.HasOutputs(2)

	thenTx.Output(1).
		HasBucket(bucket.BSV).
		HasSatoshis(change).
		UnlockableBySender()
}

func TestOutlineNoChange(t *testing.T) {
	given, then := testabilities.New(t)

	// given:
	service := given.NewTransactionOutlinesService()

	// and:
	change := bsv.Satoshis(0)
	utxoValue := bsv.Satoshis(10)

	// and:
	given.UTXOSelector().WillReturnUTXOs(change, utxoValue)

	// when:
	tx, err := service.CreateBEEF(context.Background(), given.MinimumValidTransactionSpec())

	// then:
	thenTx := then.Created(tx).WithNoError(err).WithParseableBEEFHex()

	thenTx.HasOutputs(1)
}
