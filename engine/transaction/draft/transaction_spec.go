package draft

import (
	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/draft/evaluation"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/draft/outputs"
	txerrors "github.com/bitcoin-sv/spv-wallet/engine/transaction/errors"
)

// TransactionSpec represents client provided specification for a transaction draft.
type TransactionSpec struct {
	Outputs *outputs.Specifications
}

func (t *TransactionSpec) outputs(ctx evaluation.Context) ([]*sdk.TransactionOutput, transaction.OutputsAnnotations, error) {
	if t.Outputs == nil {
		return nil, nil, txerrors.ErrDraftRequiresAtLeastOneOutput
	}

	outs, annotations, err := t.Outputs.Evaluate(ctx)
	if err != nil {
		return nil, nil, spverrors.Wrapf(err, "failed to evaluate outputs")
	}

	return outs, annotations, nil
}
