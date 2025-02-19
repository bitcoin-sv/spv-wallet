package outlines

import (
	"context"

	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/paymail"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/bsv"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/transaction"
	txerrors "github.com/bitcoin-sv/spv-wallet/engine/v2/transaction/errors"
	bsvmodel "github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/rs/zerolog"
)

type service struct {
	logger                 *zerolog.Logger
	paymailService         paymail.ServiceClient
	paymailAddressService  PaymailAddressService
	transactionBEEFService TransactionBEEFService
	utxoSelector           UTXOSelector
	feeUnit               bsvmodel.FeeUnit
	usersService          UsersService
}

// NewService creates a new transaction outlines service.
func NewService(
	paymailService paymail.ServiceClient,
	paymailAddressService PaymailAddressService,
	transactionBEEFService TransactionBEEFService,
	utxoSelector UTXOSelector,
	feeUnit bsvmodel.FeeUnit,
	logger zerolog.Logger,
	usersService UsersService,
) Service {
	if paymailService == nil {
		panic("paymail.ServiceClient is required to create transaction outlines service")
	}

	if paymailAddressService == nil {
		panic("PaymailAddressService is required to create transaction outlines service")
	}

	if transactionBEEFService == nil {
		panic("Transaction BEEF service is required to create transaction outlines service")
	}

	if utxoSelector == nil {
		panic("UTXO selector is required to create transaction outlines service")
	}

	return &service{
		logger:                 &logger,
		paymailService:         paymailService,
		paymailAddressService:  paymailAddressService,
		transactionBEEFService: transactionBEEFService,
		utxoSelector:           utxoSelector,
		feeUnit:               feeUnit,
		usersService:          usersService,
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

	beef, err := s.transactionBEEFService.PrepareBEEF(ctx, tx)
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

	evaluationCtx := s.createEvaluationContext(ctx, spec.UserID)

	tx, annotations, err := spec.evaluate(evaluationCtx)
	if err != nil {
		return nil, transaction.Annotations{}, err
	}
	return tx, annotations, err
}

func (s *service) createEvaluationContext(ctx context.Context, userID string) *evaluationContext {
	return &evaluationContext{
		Context:               ctx,
		userID:                userID,
		log:                   s.logger,
		paymail:               s.paymailService,
		paymailAddressService: s.paymailAddressService,
		utxoSelector:          s.utxoSelector,
		feeUnit:               s.feeUnit,
		usersService:          s.usersService,
	}
}
