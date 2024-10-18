package outlines

import (
	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction"
	txerrors "github.com/bitcoin-sv/spv-wallet/engine/transaction/errors"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/outlines/internal/evaluation"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/outlines/outputs"
)

// TransactionSpec represents client provided specification for a transaction outline.
type TransactionSpec struct {
	Outputs *outputs.Specifications
	XPubID  string
}

func (t *TransactionSpec) outputs(ctx evaluation.Context) ([]*sdk.TransactionOutput, transaction.OutputsAnnotations, error) {
	if t.Outputs == nil {
		return nil, nil, txerrors.ErrTxOutlineRequiresAtLeastOneOutput
	}

	outs, annotations, err := t.Outputs.Evaluate(ctx)
	if err != nil {
		return nil, nil, spverrors.Wrapf(err, "failed to evaluate outputs")
	}

	return outs, annotations, nil
}
