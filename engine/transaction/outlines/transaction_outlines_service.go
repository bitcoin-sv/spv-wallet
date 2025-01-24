package outlines

import (
	"context"

	"github.com/bitcoin-sv/spv-wallet/engine/paymail"
	"github.com/bitcoin-sv/spv-wallet/engine/paymailaddress"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	txerrors "github.com/bitcoin-sv/spv-wallet/engine/transaction/errors"
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
		panic("PaymailAddressService is required to create transaction outlines service")
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

	if spec.UserID == "" {
		return nil, txerrors.ErrTxOutlineSpecificationUserIDRequired
	}

	c := newOutlineEvaluationContext(
		ctx,
		spec.UserID,
		s.logger,
		s.paymailService,
		s.paymailAddressService,
	)

	tx, annotations, err := spec.evaluate(c)
	if err != nil {
		return nil, err
	}

	beef, err := tx.BEEFHex()
	if err != nil {
		return nil, spverrors.Wrapf(err, "failed to create transaction outline")
	}

	return &Transaction{
		BEEF:        beef,
		Annotations: annotations,
	}, nil
}
