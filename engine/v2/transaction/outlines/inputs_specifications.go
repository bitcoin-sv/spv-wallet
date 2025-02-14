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

func (s *InputsSpec) evaluate(ctx *evaluationContext, outputs annotatedOutputs) (annotatedInputs, error) {
	outs, _ := outputs.splitIntoTransactionOutputsAndAnnotations()

	tx := &sdk.Transaction{
		Outputs: outs,

		// TODO: consider creating partial transaction with only outputs makes debugging problem when debugger tries to show it by calling String() method.
		// Idea is to change UTXOSelector().Select to Select(ctx context.Context, outputsTotalValue bsv.Satoshis, byteSizeOfTxToFund uint64, userID string)
		// Size and total output satoshis can be calculated by (outputsSize(outputs) + txEnvelopeSize) outputs.totalSatoshis()
		Inputs: make([]*sdk.TransactionInput, 0),
	}

	outputs.totalSatoshis()

	utxos, err := ctx.UTXOSelector().Select(ctx, tx, ctx.UserID())
	if err != nil {
		return nil, spverrors.ErrInternal.Wrap(err)
	}

	if len(utxos) == 0 {
		return nil, txerrors.ErrTxOutlineInsufficientFunds
	}

	inputs := make(annotatedInputs, len(utxos))
	for i, utxo := range utxos {
		txID, err := chainhash.NewHashFromHex(utxo.TxID)
		if err != nil {
			return nil, spverrors.Wrapf(err, "failed to parse source transaction ID")
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

	return inputs, nil
}

type annotatedInputs []*annotatedInput

type annotatedInput struct {
	*transaction.InputAnnotation
	*sdk.TransactionInput
	utxoSatoshis  bsv.Satoshis
	estimatedSize uint64
}

func (a annotatedInputs) splitIntoTransactionInputsAndAnnotations() ([]*sdk.TransactionInput, transaction.InputAnnotations) {
	inputs := make([]*sdk.TransactionInput, len(a))
	annotationByInputIndex := make(transaction.InputAnnotations)
	for index, input := range a {
		inputs[index] = input.TransactionInput
		if input.InputAnnotation != nil {
			annotationByInputIndex[index] = input.InputAnnotation
		}
	}
	return inputs, annotationByInputIndex
}

func (a annotatedInputs) totalSatoshis() bsv.Satoshis {
	var total bsv.Satoshis
	for _, input := range a {
		total += input.utxoSatoshis
	}
	return total
}
