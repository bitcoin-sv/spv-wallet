package testabilities

import (
	"context"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures/txtestability"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/txmodels"
	"github.com/samber/lo"
	"github.com/stretchr/testify/require"
	"testing"
	"time"
)

type MockRepo struct {
	t                testing.TB
	subjectTx        txtestability.TransactionSpec
	row              *txmodels.TrackedTransaction
	fixture          *fixtureTXsync
	updated          *txmodels.TrackedTransaction
	willFailOnGet    bool
	willFailOnUpdate bool
}

func newMockRepo(t testing.TB) *MockRepo {
	return &MockRepo{
		t: t,
	}
}

func (m *MockRepo) defaultTx() txtestability.TransactionSpec {
	return txtestability.Given(m.t).Tx().WithInput(10).WithP2PKHOutput(9)
}

func (m *MockRepo) Updated() bool {
	return m.updated != nil
}

func (m *MockRepo) UpdateTransaction(ctx context.Context, trackedTx *txmodels.TrackedTransaction) error {
	if m.willFailOnUpdate {
		return spverrors.Newf("UpdateTransaction failed")
	}
	m.updated = trackedTx
	return nil
}

func (m *MockRepo) GetTransaction(_ context.Context, txID string) (transaction *txmodels.TrackedTransaction, err error) {
	if m.willFailOnGet {
		return nil, spverrors.Newf("GetTransaction failed")
	}

	require.NotNil(m.t, m.row, "Test subject transaction is not set")
	require.Equal(m.t, m.row.ID, txID, "Service asked for wrong transaction ID than expected")

	return m.row, nil
}

func (m *MockRepo) createTrackedTx() *txmodels.TrackedTransaction {
	m.subjectTx = m.defaultTx()

	return &txmodels.TrackedTransaction{
		ID:        m.subjectTx.ID(),
		TxStatus:  txmodels.TxStatusBroadcasted,
		CreatedAt: time.Now().Add(-10 * time.Minute),
		UpdatedAt: time.Now().Add(-10 * time.Minute),
	}
}

type Format string

const (
	FormatHex  Format = "hex"
	FormatBEEF Format = "beef"
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

func (m *MockRepo) WillFailOnGet() RepoFixtures {
	m.willFailOnGet = true
	return m
}

func (m *MockRepo) WillFailOnUpdate() RepoFixtures {
	m.willFailOnUpdate = true
	return m
}
