package outlines

import (
	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction"
)

// TransactionSpec represents client provided specification for a transaction outline.
type TransactionSpec struct {
	Outputs OutputsSpec
	UserID  string
	Inputs  InputsSpec
}

func (t *TransactionSpec) evaluate(ctx *evaluationContext) (*sdk.Transaction, transaction.Annotations, error) {
	outputs, err := t.Outputs.evaluate(ctx)
	if err != nil {
		return nil, transaction.Annotations{}, spverrors.Wrapf(err, "failed to evaluate outputs")
	}

	inputs, err := t.Inputs.evaluate(ctx, outputs)
	if err != nil {
		return nil, transaction.Annotations{}, err
	}

	txOuts, outputsAnnotations := outputs.splitIntoTransactionOutputsAndAnnotations()
	txIns, inputsAnnotations := inputs.splitIntoTransactionInputsAndAnnotations()

	tx := &sdk.Transaction{
		Inputs:  txIns,
		Outputs: txOuts,
	}

	annotations := transaction.Annotations{
		Inputs:  inputsAnnotations,
		Outputs: outputsAnnotations,
	}

	return tx, annotations, nil
}
