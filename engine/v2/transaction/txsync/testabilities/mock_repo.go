package testabilities

import (
	"context"
	"testing"
	"time"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures/txtestability"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/txmodels"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
)

func MockTx(t testing.TB) txtestability.TransactionSpec {
	return txtestability.Given(t).Tx().WithInput(10).WithP2PKHOutput(9)
}

type MockRepo struct {
	t          testing.TB
	subjectTx  txtestability.TransactionSpec
	row        *txmodels.TrackedTransaction
	updated    *txmodels.TrackedTransaction
	willFailOn FailingPoint
}

func newMockRepo(t testing.TB) *MockRepo {
	return &MockRepo{
		t: t,
	}
}

func (m *MockRepo) Updated() bool {
	return m.updated != nil
}

func (m *MockRepo) UpdateTransaction(ctx context.Context, trackedTx *txmodels.TrackedTransaction) error {
	if m.willFailOn == FailingPointUpdate {
		return spverrors.Newf("UpdateTransaction failed")
	}
	m.updated = trackedTx
	return nil
}

func (m *MockRepo) GetTransaction(_ context.Context, txID string) (transaction *txmodels.TrackedTransaction, err error) {
	if m.willFailOn == FailingPointGet {
		return nil, spverrors.Newf("GetTransaction failed")
	}

	require.NotNil(m.t, m.row, "Test subject transaction is not set")
	require.Equal(m.t, m.row.ID, txID, "Service asked for wrong transaction ID than expected")

	return m.row, nil
}

func (m *MockRepo) createTrackedTx() *txmodels.TrackedTransaction {
	m.subjectTx = MockTx(m.t)

	return &txmodels.TrackedTransaction{
		ID:        m.subjectTx.ID(),
		TxStatus:  txmodels.TxStatusBroadcasted,
		CreatedAt: time.Now().Add(-10 * time.Minute),
		UpdatedAt: time.Now().Add(-10 * time.Minute),
	}
}

type Format int

const (
	FormatHex Format = iota
	FormatBEEF
)

func (m *MockRepo) ContainsBroadcastedTx(format Format) RepoFixtures {
	m.row = m.createTrackedTx()
	if format == FormatHex {
		m.row.RawHex = lo.ToPtr(m.subjectTx.RawTX())
	} else {
		m.row.BeefHex = lo.ToPtr(m.subjectTx.BEEF())
	}
	return m
}

type FailingPoint int

const (
	FailingPointGet = iota + 1
	FailingPointUpdate
)

func (m *MockRepo) WillFailOn(f FailingPoint) {
	m.willFailOn = f
}
