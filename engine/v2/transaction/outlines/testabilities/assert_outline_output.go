package testabilities

import (
	"testing"

	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures/txtestability"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/bitcoin-sv/spv-wallet/models/transaction/bucket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type OutputAssertion interface {
	HasBucket(bucket bucket.Name) OutputAssertion
	HasSatoshis(satoshis bsv.Satoshis) OutputAssertion
	HasLockingScript(lockingScript string) OutputAssertion
	IsDataOnly() OutputAssertion
	IsPaymail() TransactionOutlinePaymailOutputAssertion
	UnlockableBySender() TransactionOutlinePaymailOutputAssertion
}

type TransactionOutlinePaymailOutputAssertion interface {
	HasReceiver(receiver string) TransactionOutlinePaymailOutputAssertion
	HasSender(sender string) TransactionOutlinePaymailOutputAssertion
	HasReference(reference string) TransactionOutlinePaymailOutputAssertion
}

type txOutputAssertion struct {
	t          testing.TB
	parent     *assertion
	assert     *assert.Assertions
	require    *require.Assertions
	txout      *sdk.TransactionOutput
	annotation *transaction.OutputAnnotation
	index      uint32
	txFixture  txtestability.TransactionsFixtures
}

func (a *txOutputAssertion) HasBucket(bucket bucket.Name) OutputAssertion {
	a.t.Helper()
	a.require.NotNil(a.annotation, "Output %d has no annotation", a.index)
	a.assert.Equal(bucket, a.annotation.Bucket, "Output %d has invalid bucket annotation", a.index)
	return a
}

func (a *txOutputAssertion) HasSatoshis(satoshis bsv.Satoshis) OutputAssertion {
	a.t.Helper()
	a.assert.EqualValues(satoshis, a.txout.Satoshis, "Output %d has invalid satoshis value", a.index)
	return a
}

func (a *txOutputAssertion) HasLockingScript(lockingScript string) OutputAssertion {
	a.t.Helper()
	a.assert.Equal(lockingScript, a.txout.LockingScriptHex(), "Output %d has invalid locking script", a.index)
	return a
}

func (a *txOutputAssertion) IsDataOnly() OutputAssertion {
	a.t.Helper()
	a.assert.Zerof(a.txout.Satoshis, "Output %d has value in satoshis which is not allowed for data only outputs", a.index)
	a.assert.True(a.txout.LockingScript.IsData(), "Output %d has locking script which is not data script", a.index)
	return a
}

func (a *txOutputAssertion) IsPaymail() TransactionOutlinePaymailOutputAssertion {
	a.t.Helper()
	a.require.NotNil(a.annotation, "Output %d has no annotation", a.index)
	a.require.NotNil(a.annotation.Paymail, "Output %d is not a paymail output", a.index)
	return a
}

func (a *txOutputAssertion) HasReceiver(receiver string) TransactionOutlinePaymailOutputAssertion {
	a.t.Helper()
	a.assert.Equal(receiver, a.annotation.Paymail.Receiver, "Output %d has invalid paymail receiver", a.index)
	return a
}

func (a *txOutputAssertion) HasSender(sender string) TransactionOutlinePaymailOutputAssertion {
	a.t.Helper()
	a.assert.Equal(sender, a.annotation.Paymail.Sender, "Output %d has invalid paymail sender", a.index)
	return a
}

func (a *txOutputAssertion) HasReference(reference string) TransactionOutlinePaymailOutputAssertion {
	a.t.Helper()
	a.assert.Equal(reference, a.annotation.Paymail.Reference, "Output %d has invalid paymail reference", a.index)
	return a
}

func (a *txOutputAssertion) UnlockableBySender() TransactionOutlinePaymailOutputAssertion {
	a.require.NotNil(a.annotation.CustomInstructions, "output %d has no custom instructions", a.index)

	a.txFixture.Tx().
		WithSender(fixtures.Sender).
		WithInputFromUTXO(a.parent.tx, a.index, *a.annotation.CustomInstructions...).
		WithOPReturn("dummy data").
		TX() // during TX call, the transaction is signed. Should fail if the UTXO cannot be unlocked by the user.

	return a
}
