package testabilities

import (
	"testing"

	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/draft"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/bitcoin-sv/spv-wallet/models/transaction/bucket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type DraftTransactionAssertion interface {
	Created(transaction *draft.Transaction) CreatedDraftTransactionAssertion
}

type CreatedDraftTransactionAssertion interface {
	WithNoError(err error) SuccessfullyCreatedDraftTransactionAssertion
	WithError(err error) ErrorCreationDraftTransactionAssertion
}

type ErrorCreationDraftTransactionAssertion interface {
	ThatIs(expectedError error)
}

type SuccessfullyCreatedDraftTransactionAssertion interface {
	WithParseableBEEFHex() WithParseableBEEFDraftTransactionAssertion
}

type WithParseableBEEFDraftTransactionAssertion interface {
	HasOutputs(count int) WithParseableBEEFDraftTransactionAssertion
	HasOutput(index int, assert func(OutputAssertion)) WithParseableBEEFDraftTransactionAssertion
	Output(index int) OutputAssertion
}

type OutputAssertion interface {
	HasBucket(bucket bucket.Name) OutputAssertion
	HasSatoshis(satoshis bsv.Satoshis) OutputAssertion
	HasLockingScript(lockingScript string) OutputAssertion
	IsDataOnly() OutputAssertion
	IsPaymail() DraftTransactionPaymailOutputAssertion
}

type DraftTransactionPaymailOutputAssertion interface {
	HasReceiver(receiver string) DraftTransactionPaymailOutputAssertion
	HasSender(sender string) DraftTransactionPaymailOutputAssertion
	HasReference(reference string) DraftTransactionPaymailOutputAssertion
}

func Then(t testing.TB) DraftTransactionAssertion {
	return &createdDraftAssertion{t: t, require: require.New(t), assert: assert.New(t)}
}

type createdDraftAssertion struct {
	t       testing.TB
	require *require.Assertions
	assert  *assert.Assertions
	draft   *draft.Transaction
	tx      *sdk.Transaction
	err     error
}

func (a *createdDraftAssertion) Created(transaction *draft.Transaction) CreatedDraftTransactionAssertion {
	a.draft = transaction
	return a
}

func (a *createdDraftAssertion) WithError(err error) ErrorCreationDraftTransactionAssertion {
	a.assert.Nil(a.draft)
	a.assert.Error(err)
	a.err = err
	return a
}

func (a *createdDraftAssertion) ThatIs(expectedError error) {
	a.assert.ErrorIs(a.err, expectedError)
}

// WithNoError checks if there was no error and result is not nil. It also checks if BEEF hex is parseable.
func (a *createdDraftAssertion) WithNoError(err error) SuccessfullyCreatedDraftTransactionAssertion {
	a.require.NoError(err, "Creation of draft has finished with error")
	a.require.NotNil(a.draft, "Draft should be created if there is no error")
	return a
}

func (a *createdDraftAssertion) WithParseableBEEFHex() WithParseableBEEFDraftTransactionAssertion {
	a.t.Helper()
	a.t.Logf("BEEF: %s", a.draft.BEEF)

	var err error
	a.tx, err = sdk.NewTransactionFromBEEFHex(a.draft.BEEF)
	a.require.NoErrorf(err, "Draft has invalid BEEF hex: %s", a.draft.BEEF)
	return a
}

func (a *createdDraftAssertion) HasOutputs(count int) WithParseableBEEFDraftTransactionAssertion {
	a.require.Lenf(a.tx.Outputs, count, "BEEF of draft transaction has invalid number of outputs")
	a.require.Lenf(a.draft.Annotations.Outputs, count, "Annotations of draft transaction has invalid number of outputs")
	return a
}

type draftTransactionOutputAssertion struct {
	parent     *createdDraftAssertion
	assert     *assert.Assertions
	require    *require.Assertions
	txout      *sdk.TransactionOutput
	annotation *transaction.OutputAnnotation
	index      int
}

func (a *createdDraftAssertion) HasOutput(index int, assert func(OutputAssertion)) WithParseableBEEFDraftTransactionAssertion {
	assert(a.Output(index))
	return a
}

func (a *createdDraftAssertion) Output(index int) OutputAssertion {
	a.require.Greater(len(a.tx.Outputs), index, "Draft transaction outputs has no element at index %d", index)
	a.require.Greater(len(a.draft.Annotations.Outputs), index, "Draft transaction annotation outputs has no element at index %d", index)

	return &draftTransactionOutputAssertion{
		parent:     a,
		assert:     a.assert,
		require:    a.require,
		txout:      a.tx.Outputs[index],
		annotation: a.draft.Annotations.Outputs[index],
		index:      index,
	}
}

func (a *draftTransactionOutputAssertion) HasBucket(bucket bucket.Name) OutputAssertion {
	a.assert.Equal(bucket, a.annotation.Bucket, "Output %d has invalid bucket annotation", a.index)
	return a
}

func (a *draftTransactionOutputAssertion) HasSatoshis(satoshis bsv.Satoshis) OutputAssertion {
	a.assert.EqualValues(satoshis, a.txout.Satoshis, "Output %d has invalid satoshis value", a.index)
	return a
}

func (a *draftTransactionOutputAssertion) HasLockingScript(lockingScript string) OutputAssertion {
	a.assert.Equal(lockingScript, a.txout.LockingScriptHex(), "Output %d has invalid locking script", a.index)
	return a
}

func (a *draftTransactionOutputAssertion) IsDataOnly() OutputAssertion {
	a.assert.Zerof(a.txout.Satoshis, "Output %d has value in satoshis which is not allowed for data only outputs", a.index)
	a.assert.True(a.txout.LockingScript.IsData(), "Output %d has locking script which is not data script", a.index)
	return a
}

func (a *draftTransactionOutputAssertion) IsPaymail() DraftTransactionPaymailOutputAssertion {
	a.assert.NotNil(a.annotation.Paymail, "Output %d is not a paymail output", a.index)
	return a
}

func (a *draftTransactionOutputAssertion) HasReceiver(receiver string) DraftTransactionPaymailOutputAssertion {
	if a.annotation.Paymail != nil {
		a.assert.Equal(receiver, a.annotation.Paymail.Receiver, "Output %d has invalid paymail receiver", a.index)
	}
	return a
}

func (a *draftTransactionOutputAssertion) HasSender(sender string) DraftTransactionPaymailOutputAssertion {
	if a.annotation.Paymail != nil {
		a.assert.Equal(sender, a.annotation.Paymail.Sender, "Output %d has invalid paymail sender", a.index)
	}
	return a
}

func (a *draftTransactionOutputAssertion) HasReference(reference string) DraftTransactionPaymailOutputAssertion {
	if a.annotation.Paymail != nil {
		a.assert.Equal(reference, a.annotation.Paymail.Reference, "Output %d has invalid paymail reference", a.index)
	}
	return a
}
