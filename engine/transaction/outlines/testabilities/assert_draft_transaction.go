package testabilities

import (
	"testing"

	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/outlines"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/bitcoin-sv/spv-wallet/models/transaction/bucket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type TransactionOutlineAssertion interface {
	Created(transaction *outlines.Transaction) CreatedTransactionOutlineAssertion
}

type CreatedTransactionOutlineAssertion interface {
	WithNoError(err error) SuccessfullyCreatedTransactionOutlineAssertion
	WithError(err error) ErrorCreationTransactionOutlineAssertion
}

type ErrorCreationTransactionOutlineAssertion interface {
	ThatIs(expectedError error)
}

type SuccessfullyCreatedTransactionOutlineAssertion interface {
	WithParseableBEEFHex() WithParseableBEEFTransactionOutlineAssertion
}

type WithParseableBEEFTransactionOutlineAssertion interface {
	HasOutputs(count int) WithParseableBEEFTransactionOutlineAssertion
	HasOutput(index int, assert func(OutputAssertion)) WithParseableBEEFTransactionOutlineAssertion
	Output(index int) OutputAssertion
}

type OutputAssertion interface {
	HasBucket(bucket bucket.Name) OutputAssertion
	HasSatoshis(satoshis bsv.Satoshis) OutputAssertion
	HasLockingScript(lockingScript string) OutputAssertion
	IsDataOnly() OutputAssertion
	IsPaymail() TransactionOutlinePaymailOutputAssertion
}

type TransactionOutlinePaymailOutputAssertion interface {
	HasReceiver(receiver string) TransactionOutlinePaymailOutputAssertion
	HasSender(sender string) TransactionOutlinePaymailOutputAssertion
	HasReference(reference string) TransactionOutlinePaymailOutputAssertion
}

func Then(t testing.TB) TransactionOutlineAssertion {
	return &assertion{t: t, require: require.New(t), assert: assert.New(t)}
}

type assertion struct {
	t         testing.TB
	require   *require.Assertions
	assert    *assert.Assertions
	txOutline *outlines.Transaction
	tx        *sdk.Transaction
	err       error
}

func (a *assertion) Created(transaction *outlines.Transaction) CreatedTransactionOutlineAssertion {
	a.txOutline = transaction
	return a
}

func (a *assertion) WithError(err error) ErrorCreationTransactionOutlineAssertion {
	a.assert.Nil(a.txOutline)
	a.assert.Error(err)
	a.err = err
	return a
}

func (a *assertion) ThatIs(expectedError error) {
	a.assert.ErrorIs(a.err, expectedError)
}

// WithNoError checks if there was no error and result is not nil. It also checks if BEEF hex is parseable.
func (a *assertion) WithNoError(err error) SuccessfullyCreatedTransactionOutlineAssertion {
	a.require.NoError(err, "Creation of transaction outline has finished with error")
	a.require.NotNil(a.txOutline, "Transaction outline should be created if there is no error")
	return a
}

func (a *assertion) WithParseableBEEFHex() WithParseableBEEFTransactionOutlineAssertion {
	a.t.Helper()
	a.t.Logf("BEEF: %s", a.txOutline.BEEF)

	var err error
	a.tx, err = sdk.NewTransactionFromBEEFHex(a.txOutline.BEEF)
	a.require.NoErrorf(err, "Transaction outline has invalid BEEF hex: %s", a.txOutline.BEEF)
	return a
}

func (a *assertion) HasOutputs(count int) WithParseableBEEFTransactionOutlineAssertion {
	a.require.Lenf(a.tx.Outputs, count, "BEEF of transaction outline has invalid number of outputs")
	a.require.Lenf(a.txOutline.Annotations.Outputs, count, "Annotations of transaction outline has invalid number of outputs")
	return a
}

type txOutputAssertion struct {
	parent     *assertion
	assert     *assert.Assertions
	require    *require.Assertions
	txout      *sdk.TransactionOutput
	annotation *transaction.OutputAnnotation
	index      int
}

func (a *assertion) HasOutput(index int, assert func(OutputAssertion)) WithParseableBEEFTransactionOutlineAssertion {
	assert(a.Output(index))
	return a
}

func (a *assertion) Output(index int) OutputAssertion {
	a.require.Greater(len(a.tx.Outputs), index, "Transaction outline outputs has no element at index %d", index)
	a.require.Greater(len(a.txOutline.Annotations.Outputs), index, "Transaction outline annotation outputs has no element at index %d", index)

	return &txOutputAssertion{
		parent:     a,
		assert:     a.assert,
		require:    a.require,
		txout:      a.tx.Outputs[index],
		annotation: a.txOutline.Annotations.Outputs[index],
		index:      index,
	}
}

func (a *txOutputAssertion) HasBucket(bucket bucket.Name) OutputAssertion {
	a.assert.Equal(bucket, a.annotation.Bucket, "Output %d has invalid bucket annotation", a.index)
	return a
}

func (a *txOutputAssertion) HasSatoshis(satoshis bsv.Satoshis) OutputAssertion {
	a.assert.EqualValues(satoshis, a.txout.Satoshis, "Output %d has invalid satoshis value", a.index)
	return a
}

func (a *txOutputAssertion) HasLockingScript(lockingScript string) OutputAssertion {
	a.assert.Equal(lockingScript, a.txout.LockingScriptHex(), "Output %d has invalid locking script", a.index)
	return a
}

func (a *txOutputAssertion) IsDataOnly() OutputAssertion {
	a.assert.Zerof(a.txout.Satoshis, "Output %d has value in satoshis which is not allowed for data only outputs", a.index)
	a.assert.True(a.txout.LockingScript.IsData(), "Output %d has locking script which is not data script", a.index)
	return a
}

func (a *txOutputAssertion) IsPaymail() TransactionOutlinePaymailOutputAssertion {
	a.assert.NotNil(a.annotation.Paymail, "Output %d is not a paymail output", a.index)
	return a
}

func (a *txOutputAssertion) HasReceiver(receiver string) TransactionOutlinePaymailOutputAssertion {
	if a.annotation.Paymail != nil {
		a.assert.Equal(receiver, a.annotation.Paymail.Receiver, "Output %d has invalid paymail receiver", a.index)
	}
	return a
}

func (a *txOutputAssertion) HasSender(sender string) TransactionOutlinePaymailOutputAssertion {
	if a.annotation.Paymail != nil {
		a.assert.Equal(sender, a.annotation.Paymail.Sender, "Output %d has invalid paymail sender", a.index)
	}
	return a
}

func (a *txOutputAssertion) HasReference(reference string) TransactionOutlinePaymailOutputAssertion {
	if a.annotation.Paymail != nil {
		a.assert.Equal(reference, a.annotation.Paymail.Reference, "Output %d has invalid paymail reference", a.index)
	}
	return a
}
