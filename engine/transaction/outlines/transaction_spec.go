package outlines

import (
	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction"
)

// TransactionSpec represents client provided specification for a transaction outline.
type TransactionSpec struct {
	Outputs OutputsSpec
	UserID  string
}

func (t *TransactionSpec) evaluate(ctx evaluationContext) (*sdk.Transaction, transaction.Annotations, error) {
	outputs, err := t.Outputs.evaluate(ctx)
	if err != nil {
		return nil, transaction.Annotations{}, spverrors.Wrapf(err, "failed to evaluate outputs")
	}

	txOuts, outputsAnnotations := outputs.splitIntoTransactionOutputsAndAnnotations()
	tx := &sdk.Transaction{
		Outputs: txOuts,
	}

	annotations := transaction.Annotations{
		Outputs: outputsAnnotations,
	}

	return tx, annotations, nil
}
