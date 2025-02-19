package outlines

import (
	"github.com/bitcoin-sv/go-sdk/chainhash"
	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction"
	txerrors "github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/errors"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
)

// InputsSpec are representing a client specification for inputs part of the transaction.
type InputsSpec struct {
}

func (s *InputsSpec) evaluate(ctx *evaluationContext, outputs annotatedOutputs) (annotatedInputs, bsv.Satoshis, error) {
	outs := outputs.toTransactionOutputs()

	tx := &sdk.Transaction{
		Outputs: outs,
		Inputs:  make([]*sdk.TransactionInput, 0),
	}

	utxos, change, err := ctx.UTXOSelector().Select(ctx, tx, ctx.UserID())
	if err != nil {
		return nil, 0, spverrors.ErrInternal.Wrap(err)
	}

	if len(utxos) == 0 {
		return nil, 0, txerrors.ErrTxOutlineInsufficientFunds
	}

	inputs := make(annotatedInputs, len(utxos))
	for i, utxo := range utxos {
		txID, err := chainhash.NewHashFromHex(utxo.TxID)
		if err != nil {
			return nil, 0, spverrors.Wrapf(err, "failed to parse source transaction ID")
		}
		inputs[i] = &annotatedInput{
			TransactionInput: &sdk.TransactionInput{
				SourceTXID:       txID,
				SourceTxOutIndex: utxo.Vout,
			},
			InputAnnotation: &transaction.InputAnnotation{
				CustomInstructions: utxo.CustomInstructions,
			},
			utxoSatoshis:  utxo.Satoshis,
			estimatedSize: utxo.EstimatedInputSize,
		}
	}

	return inputs, change, nil
}

type annotatedInputs []*annotatedInput

type annotatedInput struct {
	*transaction.InputAnnotation
	*sdk.TransactionInput
	utxoSatoshis  bsv.Satoshis
	estimatedSize uint64
}

func (a annotatedInputs) splitIntoTransactionInputsAndAnnotations() ([]*sdk.TransactionInput, transaction.InputAnnotations) {
	return a.toTransactionInputs(), a.toAnnotations()
}

func (a annotatedInputs) toTransactionInputs() []*sdk.TransactionInput {
	inputs := make([]*sdk.TransactionInput, len(a))
	for i, in := range a {
		inputs[i] = in.TransactionInput
	}
	return inputs
}

func (a annotatedInputs) toAnnotations() transaction.InputAnnotations {
	annotations := make(transaction.InputAnnotations)
	for inputIndex, in := range a {
		if in.InputAnnotation != nil {
			annotations[inputIndex] = in.InputAnnotation
		}
	}
	return annotations
}
