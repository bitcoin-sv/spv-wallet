package testabilities

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/engine/database"
	"github.com/stretchr/testify/require"
)

type ErrorAssert interface {
	NothingChanged()
}

type RecordOutlineAssert interface {
	WithNoError(err error) SuccessfullyCreatedRecordOutlineAssertion
	WithErrorIs(err, expectedError error) ErrorAssert

	StoredOutputs([]database.Output) RecordOutlineAssert
	StoredData([]database.Data) RecordOutlineAssert
}

type SuccessfullyCreatedRecordOutlineAssertion interface {
	Broadcasted(txID string) SuccessfullyCreatedRecordOutlineAssertion
	StoredAsBroadcasted(txID string) SuccessfullyCreatedRecordOutlineAssertion
}

type assert struct {
	t       testing.TB
	require *require.Assertions
	given   *recordServiceFixture
}

func then(t testing.TB, given *recordServiceFixture) RecordOutlineAssert {
	return &assert{
		t:       t,
		require: require.New(t),
		given:   given,
	}
}

func (a *assert) WithNoError(err error) SuccessfullyCreatedRecordOutlineAssertion {
	a.require.NoError(err, "Record transaction outline has error")
	return a
}

func (a *assert) WithErrorIs(err, expectedError error) ErrorAssert {
	require.Error(a.t, err, "Record transaction outline has no error")
	require.ErrorIs(a.t, err, expectedError, "Record transaction outline has wrong error")
	return a
}

func (a *assert) Broadcasted(txID string) SuccessfullyCreatedRecordOutlineAssertion {
	tx := a.given.broadcaster.checkBroadcasted(txID)
	require.NotNil(a.t, tx, "Transaction %s is not broadcasted", txID)
	return a
}

func (a *assert) StoredAsBroadcasted(txID string) SuccessfullyCreatedRecordOutlineAssertion {
	tx := a.given.repository.getTransaction(txID)
	require.NotNil(a.t, tx, "Transaction %s is not stored", txID)
	require.Equal(a.t, txID, tx.ID, "Transaction %s has wrong ID", txID)
	require.Equal(a.t, database.TxStatusBroadcasted, tx.TxStatus, "Transaction %s is not stored as broadcasted", txID)
	return a
}

func (a *assert) StoredOutputs(outputs []database.Output) RecordOutlineAssert {
	require.Subset(a.t, a.given.repository.GetAllOutputs(), outputs)
	return a
}

func (a *assert) StoredData(data []database.Data) RecordOutlineAssert {
	require.Subset(a.t, a.given.repository.GetAllData(), data)
	return a
}

func (a *assert) NothingChanged() {
	require.ElementsMatch(a.t, a.given.initialOutputs, a.given.repository.GetAllOutputs(), "Outputs are changed")
	require.ElementsMatch(a.t, a.given.initialData, a.given.repository.GetAllData(), "Data are changed")
}
