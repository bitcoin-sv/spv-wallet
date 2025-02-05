package outlines

import (
	"context"

	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/paymail"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/bsv"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/errors"
	"github.com/rs/zerolog"
)

type service struct {
	logger                *zerolog.Logger
	paymailService        paymail.ServiceClient
	paymailAddressService PaymailAddressService
	utxoSelector          UTXOSelector
}

// NewService creates a new transaction outlines service.
func NewService(paymailService paymail.ServiceClient, paymailAddressService PaymailAddressService, utxoSelector UTXOSelector, logger zerolog.Logger) Service {
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
		utxoSelector:          utxoSelector,
	}
}

func (s *service) CreateRawTx(ctx context.Context, spec *TransactionSpec) (*Transaction, error) {
	tx, annotations, err := s.evaluateSpec(ctx, spec)
	if err != nil {
		return nil, err
	}

	return &Transaction{
		Hex:         bsv.TxHex(tx.Hex()),
		Annotations: annotations,
	}, nil
}

// CreateBEEF creates a new transaction outline based on specification.
func (s *service) CreateBEEF(ctx context.Context, spec *TransactionSpec) (*Transaction, error) {
	tx, annotations, err := s.evaluateSpec(ctx, spec)
	if err != nil {
		return nil, err
	}

	beef, err := s.formatAsBEEF(tx)
	if err != nil {
		return nil, spverrors.Wrapf(err, "failed to make BEEF format for transaction outline")
	}

	return &Transaction{
		Hex:         bsv.TxHex(beef),
		Annotations: annotations,
	}, nil
}

func (s *service) evaluateSpec(ctx context.Context, spec *TransactionSpec) (*sdk.Transaction, transaction.Annotations, error) {
	if spec == nil {
		return nil, transaction.Annotations{}, txerrors.ErrTxOutlineSpecificationRequired
	}

	if spec.UserID == "" {
		return nil, transaction.Annotations{}, txerrors.ErrTxOutlineSpecificationUserIDRequired
	}

	c := newOutlineEvaluationContext(
		ctx,
		spec.UserID,
		s.logger,
		s.paymailService,
		s.paymailAddressService,
		s.utxoSelector,
	)

	tx, annotations, err := spec.evaluate(c)
	if err != nil {
		return nil, transaction.Annotations{}, err
	}
	return tx, annotations, err
}

func (s *service) formatAsBEEF(tx *sdk.Transaction) (string, error) {
	// FIXME: [waiting for SPV-1370] temporary solution to not fail on BEEF until we have service for collecting ancestors for inputs.
	tmpTx := &sdk.Transaction{
		Outputs: tx.Outputs,
	}
	return tmpTx.BEEFHex() //nolint:wrapcheck // temporary solution - will be removed after SPV-1370
}
