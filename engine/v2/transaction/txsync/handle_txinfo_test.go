package txsync_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	chainmodels "github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/txmodels"
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
			spec := testabilities.MinedTXInfo(t)

			// when:
			err := service.Handle(context.Background(), chainmodels.TXInfo(spec))

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
			spec := testabilities.TXInfo(t, status)

			// when:
			err := service.Handle(context.Background(), chainmodels.TXInfo(spec))

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
		t.Run(fmt.Sprintf("update as problematic for %v status", status), func(t *testing.T) {
			given, then := testabilities.New(t)
			// given:
			service := given.Service()

			// and:
			given.Repo().ContainsBroadcastedTx(testabilities.FormatHex)

			// and:
			spec := testabilities.TXInfo(t, status)

			// when:
			err := service.Handle(context.Background(), chainmodels.TXInfo(spec))

			// then:
			then.WithNoError(err).
				TransactionUpdated(txmodels.TxStatusProblematic)
		})
	}
}

func TestDbFails(t *testing.T) {
	tests := map[string]struct {
		spec         testabilities.TXInfoSpec
		failingPoint testabilities.FailingPoint
	}{
		"fails on get tx from DB": {
			spec:         testabilities.MinedTXInfo(t),
			failingPoint: testabilities.FailingPointGet,
		},
		"fails on get tx from DB for problematic tx": {
			spec:         testabilities.TXInfo(t, chainmodels.DoubleSpendAttempted),
			failingPoint: testabilities.FailingPointGet,
		},
		"fails on update tx in DB": {
			spec:         testabilities.MinedTXInfo(t),
			failingPoint: testabilities.FailingPointUpdate,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			given, then := testabilities.New(t)
			// given:
			service := given.Service()

			// and:
			given.Repo().
				ContainsBroadcastedTx(testabilities.FormatBEEF).
				WillFailOn(test.failingPoint)

			// when:
			err := service.Handle(context.Background(), chainmodels.TXInfo(test.spec))

			// then:
			then.WithError(err)
		})
	}
}

func TestForWrongTxInfo(t *testing.T) {
	tests := map[string]struct {
		spec testabilities.TXInfoSpec
	}{
		"empty tx info": {
			spec: testabilities.EmptyTXInfo(),
		},
		"unequal txInfo block height and the one in BUMP": {
			spec: testabilities.MinedTXInfo(t).IncrementBlockHeight(),
		},
		"wrong merkle path": {
			spec: testabilities.MinedTXInfo(t).WithWrongMerklePath(),
		},
		"subject tx out of merkle path": {
			spec: testabilities.MinedTXInfo(t).WithSubjectTxOutOfMerklePath(),
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			given, then := testabilities.New(t)
			// given:
			service := given.Service()

			// and:
			given.Repo().ContainsBroadcastedTx(testabilities.FormatBEEF)

			// when:
			err := service.Handle(context.Background(), chainmodels.TXInfo(test.spec))

			// then:
			then.WithError(err)
		})
	}
}

func TestDoNothingOnCallbackWithOldStatus(t *testing.T) {
	given, then := testabilities.New(t)
	// given:
	service := given.Service()

	// and:
	given.Repo().ContainsBroadcastedTx(testabilities.FormatBEEF)

	// and:
	spec := testabilities.MinedTXInfo(t).WithTimestamp(time.Now().Add(-1 * time.Hour))

	// when:
	err := service.Handle(context.Background(), chainmodels.TXInfo(spec))

	// then:
	then.WithNoError(err).TransactionNotUpdated()
}
