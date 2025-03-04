package txsync_test

import (
	"context"
	"fmt"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/txmodels"
	"testing"
	"time"

	chainmodels "github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/txsync/testabilities"
)

func TestUpdateOnMinedTx(t *testing.T) {
	tests := map[string]struct {
		format testabilities.Format
	}{
		"BEEF": {
			format: testabilities.FormatBEEF,
		},
		"RawHex": {
			format: testabilities.FormatHex,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			given, then := testabilities.New(t)
			// given:
			service := given.Service()

			// and:
			given.Repo().ContainsBroadcastedTx(test.format)

			// and:
			txInfo := given.MinedTXInfo()

			// when:
			err := service.Handle(context.Background(), chainmodels.TXInfo(txInfo))

			// then:
			then.WithNoError(err).
				TransactionUpdated(txmodels.TxStatusMined).
				HasBlockHash().
				HasBlockHeight().
				HasBEEF().
				HasEmptyRawHex()
		})
	}
}

func TestHandleNotUpdateOnNeutralStatuses(t *testing.T) {
	statuses := []chainmodels.TXStatus{
		chainmodels.SeenOnNetwork,
		chainmodels.Queued,
		chainmodels.Received,
		chainmodels.Stored,
		chainmodels.AnnouncedToNetwork,
		chainmodels.RequestedByNetwork,
		chainmodels.SentToNetwork,
		chainmodels.AcceptedByNetwork,
	}
	for _, status := range statuses {
		t.Run(fmt.Sprintf("not update on %v status", status), func(t *testing.T) {
			given, then := testabilities.New(t)
			// given:
			service := given.Service()

			// and:
			given.Repo().ContainsBroadcastedTx(testabilities.FormatBEEF)

			// and:
			txInfo := given.TXInfo(status)

			// when:
			err := service.Handle(context.Background(), chainmodels.TXInfo(txInfo))

			// then:
			then.WithNoError(err).
				TransactionNotUpdated()
		})
	}
}

func TestHandleNotUpdateOnProblematicStatuses(t *testing.T) {
	statuses := []chainmodels.TXStatus{
		chainmodels.Rejected,
		chainmodels.DoubleSpendAttempted,
		chainmodels.Unknown,
		chainmodels.SeenInOrphanMempool,
	}
	for _, status := range statuses {
		t.Run(fmt.Sprintf("not update on %v status", status), func(t *testing.T) {
			given, then := testabilities.New(t)
			// given:
			service := given.Service()

			// and:
			given.Repo().ContainsBroadcastedTx(testabilities.FormatHex)

			// and:
			txInfo := given.TXInfo(status)

			// when:
			err := service.Handle(context.Background(), chainmodels.TXInfo(txInfo))

			// then:
			then.WithNoError(err).
				TransactionUpdated(txmodels.TxStatusProblematic)
		})
	}
}

func TestEmptyTxID(t *testing.T) {
	given, then := testabilities.New(t)
	// given:
	service := given.Service()

	// and:
	given.Repo().ContainsBroadcastedTx(testabilities.FormatBEEF)

	// and:
	txInfo := given.EmptyTXInfo()

	// when:
	err := service.Handle(context.Background(), chainmodels.TXInfo(txInfo))

	// then:
	then.WithError(err)
}

func TestForCallbackWithOldStatus(t *testing.T) {
	given, then := testabilities.New(t)
	// given:
	service := given.Service()

	// and:
	given.Repo().ContainsBroadcastedTx(testabilities.FormatBEEF)

	// and:
	txInfo := given.MinedTXInfo().
		WithTimestamp(time.Now().Add(-1 * time.Hour))

	// when:
	err := service.Handle(context.Background(), chainmodels.TXInfo(txInfo))

	// then:
	then.WithNoError(err).
		TransactionNotUpdated()
}

func TestForUnequalBlockHeights(t *testing.T) {
	given, then := testabilities.New(t)
	// given:
	service := given.Service()

	// and:
	given.Repo().ContainsBroadcastedTx(testabilities.FormatBEEF)

	// and:
	txInfo := given.MinedTXInfo()
	txInfo.BlockHeight++

	// when:
	err := service.Handle(context.Background(), chainmodels.TXInfo(txInfo))

	// then:
	then.WithError(err)
}

func TestDbFailsOnGet(t *testing.T) {
	given, then := testabilities.New(t)
	// given:
	service := given.Service()

	// and:
	given.Repo().
		ContainsBroadcastedTx(testabilities.FormatBEEF).
		WillFailOnGet()

	// and:
	txInfo := given.MinedTXInfo()

	// when:
	err := service.Handle(context.Background(), chainmodels.TXInfo(txInfo))

	// then:
	then.WithError(err)
}

func TestDbFailsOnUpdateForMinedTx(t *testing.T) {
	given, then := testabilities.New(t)
	// given:
	service := given.Service()

	// and:
	given.Repo().
		ContainsBroadcastedTx(testabilities.FormatBEEF).
		WillFailOnUpdate()

	// and:
	txInfo := given.MinedTXInfo()

	// when:
	err := service.Handle(context.Background(), chainmodels.TXInfo(txInfo))

	// then:
	then.WithError(err)
}

func TestDbFailsOnUpdateForProblematicTx(t *testing.T) {
	given, then := testabilities.New(t)
	// given:
	service := given.Service()

	// and:
	given.Repo().
		ContainsBroadcastedTx(testabilities.FormatBEEF).
		WillFailOnUpdate()

	// and:
	txInfo := given.TXInfo(chainmodels.DoubleSpendAttempted)

	// when:
	err := service.Handle(context.Background(), chainmodels.TXInfo(txInfo))

	// then:
	then.WithError(err)
}

func TestForWrongMerklePath(t *testing.T) {
	given, then := testabilities.New(t)
	// given:
	service := given.Service()

	// and:
	given.Repo().ContainsBroadcastedTx(testabilities.FormatBEEF)

	// and:
	txInfo := given.MinedTXInfo().WithWrongMerklePath()

	// when:
	err := service.Handle(context.Background(), chainmodels.TXInfo(txInfo))

	// then:
	then.WithError(err)
}

func TestForSubjectTxOutOfMerklePath(t *testing.T) {
	given, then := testabilities.New(t)
	// given:
	service := given.Service()

	// and:
	given.Repo().ContainsBroadcastedTx(testabilities.FormatBEEF)

	// and:
	txInfo := given.MinedTXInfo().WithSubjectTxOutOfMerklePath()

	// when:
	err := service.Handle(context.Background(), chainmodels.TXInfo(txInfo))

	// then:
	then.WithError(err)
}
