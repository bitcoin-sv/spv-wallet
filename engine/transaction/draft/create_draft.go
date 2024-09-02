package draft

import (
	"context"

	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction"
	txerrors "github.com/bitcoin-sv/spv-wallet/engine/transaction/errors"
)

// Create creates a new draft transaction based on specification.
func Create(ctx context.Context, spec *TransactionSpec) (*Transaction, error) {
	tx := &sdk.Transaction{}
	if spec == nil {
		return nil, txerrors.ErrDraftSpecificationRequired
	}
	outputs, annotations, err := spec.outputs(ctx)
	if err != nil {
		return nil, err
	}
	tx.Outputs = outputs

	beef, err := tx.BEEFHex()
	if err != nil {
		return nil, spverrors.Wrapf(err, "failed to create draft transaction")
	}

	return &Transaction{
		BEEF: beef,
		Annotations: &transaction.Annotations{
			Outputs: annotations,
		},
	}, nil
}
