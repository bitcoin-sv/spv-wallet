package draft

import (
	"context"

	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/paymail"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/draft/evaluation"
	txerrors "github.com/bitcoin-sv/spv-wallet/engine/transaction/errors"
	"github.com/rs/zerolog"
)

type service struct {
	logger         *zerolog.Logger
	paymailService paymail.ServiceClient
}

// NewDraftService creates a new draft service.
func NewDraftService(paymailService paymail.ServiceClient, logger zerolog.Logger) Service {
	if paymailService == nil {
		panic("paymail service is required to create draft transaction service")
	}

	return &service{
		logger:         &logger,
		paymailService: paymailService,
	}
}

// Create creates a new draft transaction based on specification.
func (s *service) Create(ctx context.Context, spec *TransactionSpec) (*Transaction, error) {
	if spec == nil {
		return nil, txerrors.ErrDraftSpecificationRequired
	}

	c := evaluation.NewContext(ctx, s.logger, s.paymailService)

	outputs, annotations, err := spec.outputs(c)
	if err != nil {
		return nil, err
	}

	tx := &sdk.Transaction{
		Outputs: outputs,
	}

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
