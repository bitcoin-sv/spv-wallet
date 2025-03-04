package testabilities

import (
	trx "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/txmodels"
	"github.com/stretchr/testify/require"
	"testing"
)

type AssertTXsync interface {
	WithNoError(err error) AssertSucceededTXsync
	WithError(err error)
}

type AssertSucceededTXsync interface {
	TransactionUpdated(expectedStatus txmodels.TxStatus) AssertUpdatedTX
	TransactionNotUpdated()
}

type AssertUpdatedTX interface {
	HasBlockHash() AssertUpdatedTX
	HasBlockHeight() AssertUpdatedTX
	HasBEEF() AssertUpdatedTX
	HasEmptyRawHex() AssertUpdatedTX
}

func Then(t testing.TB, given FixtureTXsync) AssertTXsync {
	return &assertTXsync{
		t:       t,
		require: require.New(t),
		given:   given.(*fixtureTXsync),
	}
}

type assertTXsync struct {
	t       testing.TB
	require *require.Assertions
	given   *fixtureTXsync
}

func (a *assertTXsync) WithError(err error) {
	require.Error(a.t, err)
}

func (a *assertTXsync) WithNoError(err error) AssertSucceededTXsync {
	a.require.NoError(err)
	return a
}

func (a *assertTXsync) TransactionUpdated(expectedStatus txmodels.TxStatus) AssertUpdatedTX {
	a.require.NotNil(a.given.repo.updated, "Transaction not updated")
	updated := *a.given.repo.updated
	a.require.Equal(a.given.repo.subjectTx.ID(), updated.ID)
	a.require.Equal(expectedStatus, updated.TxStatus)
	return a
}

func (a *assertTXsync) TransactionNotUpdated() {
	a.require.False(a.given.repo.Updated())
}

func (a *assertTXsync) HasBlockHash() AssertUpdatedTX {
	a.require.NotNil(a.given.repo.updated.BlockHash)
	a.require.Equal(mockBlockHash, *a.given.repo.updated.BlockHash)
	return a
}

func (a *assertTXsync) HasBlockHeight() AssertUpdatedTX {
	a.require.NotNil(a.given.repo.updated.BlockHeight)
	a.require.Equal(int64(mockBlockHeight), *a.given.repo.updated.BlockHeight)
	return a
}

func (a *assertTXsync) HasBEEF() AssertUpdatedTX {
	a.require.NotNil(a.given.repo.updated.BeefHex)

	tx, err := trx.NewTransactionFromBEEFHex(*a.given.repo.updated.BeefHex)
	a.require.NoError(err)

	a.require.Equal(a.given.repo.subjectTx.RawTX(), tx.Hex())
	a.require.NotNil(tx.MerklePath)
	a.require.Equal(mockBump(a.given.repo.subjectTx.ID()).Hex(), tx.MerklePath.Hex())
	return a
}

func (a *assertTXsync) HasEmptyRawHex() AssertUpdatedTX {
	a.require.Nil(a.given.repo.updated.RawHex)
	return a
}
