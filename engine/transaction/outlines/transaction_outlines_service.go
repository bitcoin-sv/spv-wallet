package outlines

import (
	"context"

	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/paymail"
	"github.com/bitcoin-sv/spv-wallet/engine/paymailaddress"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction"
	txerrors "github.com/bitcoin-sv/spv-wallet/engine/transaction/errors"
	"github.com/bitcoin-sv/spv-wallet/engine/transaction/outlines/internal/evaluation"
	"github.com/rs/zerolog"
)

type service struct {
	logger                *zerolog.Logger
	paymailService        paymail.ServiceClient
	paymailAddressService paymailaddress.Service
}

// NewService creates a new transaction outlines service.
func NewService(paymailService paymail.ServiceClient, paymailAddressService paymailaddress.Service, logger zerolog.Logger) Service {
	if paymailService == nil {
		panic("paymail.ServiceClient is required to create transaction outlines service")
	}

	if paymailAddressService == nil {
		panic("paymailaddress.Service is required to create transaction outlines service")
	}

	return &service{
		logger:                &logger,
		paymailService:        paymailService,
		paymailAddressService: paymailAddressService,
	}
}

// Create creates a new transaction outline based on specification.
func (s *service) Create(ctx context.Context, spec *TransactionSpec) (*Transaction, error) {
	if spec == nil {
		return nil, txerrors.ErrTxOutlineSpecificationRequired
	}

	if spec.XPubID == "" {
		return nil, txerrors.ErrTxOutlineSpecificationXPubIDRequired
	}

	c := evaluation.NewContext(
		ctx,
		spec.XPubID,
		s.logger,
		s.paymailService,
		s.paymailAddressService,
	)

	outputs, annotations, err := spec.outputs(c)
	if err != nil {
		return nil, err
	}

	tx := &sdk.Transaction{
		Outputs: outputs,
	}

	beef, err := tx.BEEFHex()
	if err != nil {
		return nil, spverrors.Wrapf(err, "failed to create transaction outline")
	}

	return &Transaction{
		BEEF: beef,
		Annotations: transaction.Annotations{
			Outputs: annotations,
		},
	}, nil
}
