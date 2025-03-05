package paymailserver

import (
	"context"
	"fmt"

	paymailserver "github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/go-paymail/server"
	"github.com/bitcoin-sv/go-paymail/spv"
	primitives "github.com/bitcoin-sv/go-sdk/primitives/ec"
	"github.com/bitcoin-sv/go-sdk/script"
	trx "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/go-sdk/transaction/template/p2pkh"
	"github.com/bitcoin-sv/spv-wallet/engine/paymail"
	pmerrors "github.com/bitcoin-sv/spv-wallet/engine/paymail/errors"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/addresses/addressesmodels"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/keys/type42"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/paymails/paymailsmodels"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/rs/zerolog"
)

// ServiceProvider is a service provider for paymail service
type ServiceProvider interface {
	server.PaymailServiceProvider
	server.PikeContactServiceProvider
}

// NewServiceProvider create a new paymail service server which handlers incoming paymail requests
func NewServiceProvider(
	logger *zerolog.Logger,
	paymails paymail.PaymailsService,
	users paymail.UsersService,
	addresses paymail.AddressesService,
	contacts paymail.ContactsService,
	spv paymail.MerkleRootsVerifier,
	recorder paymail.TxRecorder,
) ServiceProvider {
	return &serviceProvider{
		logger:    logger,
		paymails:  paymails,
		users:     users,
		addresses: addresses,
		contacts:  contacts,
		spv:       spv,
		recorder:  recorder,
	}
}

type serviceProvider struct {
	logger    *zerolog.Logger
	paymails  paymail.PaymailsService
	users     paymail.UsersService
	addresses paymail.AddressesService
	contacts  paymail.ContactsService
	spv       paymail.MerkleRootsVerifier
	recorder  paymail.TxRecorder
}

func (s *serviceProvider) CreateAddressResolutionResponse(ctx context.Context, alias, domain string, _ bool, _ *server.RequestMetadata) (*paymailserver.ResolutionPayload, error) {
	destination, err := s.createDestinationForUser(ctx, alias, domain)
	if err != nil {
		return nil, err
	}

	return &paymailserver.ResolutionPayload{
		Address:   destination.address,
		Output:    destination.lockingScript,
		Signature: "", // signature is not supported due to "noncustodial" nature of the wallet; private keys are not stored
	}, nil
}

func (s *serviceProvider) CreateP2PDestinationResponse(ctx context.Context, alias, domain string, satoshis uint64, _ *server.RequestMetadata) (*paymailserver.PaymentDestinationPayload, error) {
	destination, err := s.createDestinationForUser(ctx, alias, domain)
	if err != nil {
		return nil, err
	}

	return &paymailserver.PaymentDestinationPayload{
		Outputs: []*paymailserver.PaymentOutput{{
			Address:  destination.address,
			Satoshis: satoshis,
			Script:   destination.lockingScript,
		}},
		Reference: destination.referenceID,
	}, nil
}

func (s *serviceProvider) GetPaymailByAlias(ctx context.Context, alias, domain string, _ *server.RequestMetadata) (*paymailserver.AddressInformation, error) {
	model, err := s.paymails.Find(ctx, alias, domain)
	if err != nil {
		return nil, pmerrors.ErrPaymailDBFailed.Wrap(err)
	}
	if model == nil {
		return nil, pmerrors.ErrPaymailNotFound
	}

	pki, _, err := s.pki(ctx, model)
	if err != nil {
		return nil, err
	}

	return &paymailserver.AddressInformation{
		Alias:  model.Alias,
		Avatar: model.Avatar,
		Domain: model.Domain,
		ID:     fmt.Sprintf("%d", model.ID),
		Name:   model.PublicName,
		PubKey: pki.ToDERHex(),
	}, nil
}

func (s *serviceProvider) RecordTransaction(ctx context.Context, p2pTx *paymailserver.P2PTransaction, requestMetadata *server.RequestMetadata) (*paymailserver.P2PTransactionPayload, error) {
	isBEEF := p2pTx.DecodedBeef != nil && p2pTx.Beef != ""
	isRawTX := p2pTx.Hex != ""

	if !isBEEF && !isRawTX {
		return nil, pmerrors.ErrParseIncomingTransaction
	}

	var tx *trx.Transaction
	var err error
	if isBEEF {
		tx, err = trx.NewTransactionFromBEEFHex(p2pTx.Beef)
	} else {
		panic("not implemented yet")
	}

	if err != nil {
		return nil, pmerrors.ErrParseIncomingTransaction.Wrap(err)
	}

	receiverPaymail := requestMetadata.Alias + "@" + requestMetadata.Domain

	err = s.recorder.RecordPaymailTransaction(ctx, tx, p2pTx.MetaData.Sender, receiverPaymail)
	if err != nil {
		return nil, pmerrors.ErrRecordTransaction.Wrap(err)
	}

	// TODO: TrackMissingTxs for BEEF purposes (or handle it in other way)

	return &paymailserver.P2PTransactionPayload{
		Note: p2pTx.MetaData.Note,
		TxID: tx.TxID().String(),
	}, nil
}

func (s *serviceProvider) VerifyMerkleRoots(ctx context.Context, merkleProofs []*spv.MerkleRootConfirmationRequestItem) error {
	// TODO include metrics for VerifyMerkleRoots (perhaps on another level - maybe ChainService)

	valid, err := s.spv.VerifyMerkleRoots(ctx, merkleProofs)

	// NOTE: these errors goes to go-paymail and are not logged there, so we need to log them here

	if err != nil {
		s.logger.Error().Err(err).Msg("Error verifying merkle roots")
		return pmerrors.ErrPaymailMerkleRootVerificationFailed.Wrap(err)
	}

	if !valid {
		s.logger.Warn().Msg("Not all provided merkle roots were confirmed by BHS")
		return pmerrors.ErrPaymailInvalidMerkleRoots
	}

	return nil
}

// AddContact is a method to add a new contact to the pike contact list
func (s *serviceProvider) AddContact(ctx context.Context, requesterPaymail string, contact *paymailserver.PikeContactRequestPayload) error {
	rAlias, rDomain, _ := paymailserver.SanitizePaymail(requesterPaymail)
	pAddress, err := s.paymails.Find(ctx, rAlias, rDomain)
	if err != nil || pAddress == nil {
		return spverrors.ErrCouldNotFindPaymail
	}

	if _, err = s.contacts.AddContactRequest(ctx, contact.FullName, contact.Paymail, pAddress.UserID); err != nil {
		return spverrors.ErrAddingContactRequest.WithTrace(err)
	}
	return nil
}

func (s *serviceProvider) pki(ctx context.Context, paymailModel *paymailsmodels.Paymail) (*primitives.PublicKey, string, error) {
	userPubKey, err := s.users.GetPubKey(ctx, paymailModel.UserID)
	if err != nil {
		return nil, "", pmerrors.ErrPaymailPKI.Wrap(err)
	}

	pki, derivationKey, err := type42.PaymailPKI(userPubKey, paymailModel.Alias, paymailModel.Domain)
	if err != nil {
		return nil, derivationKey, pmerrors.ErrPaymailPKI.Wrap(err)
	}
	return pki, derivationKey, nil
}

type destinationData struct {
	address       string
	lockingScript string
	referenceID   string
}

func (s *serviceProvider) createDestinationForUser(ctx context.Context, alias, domain string) (*destinationData, error) {
	paymailModel, err := s.paymails.Find(ctx, alias, domain)
	if err != nil {
		return nil, pmerrors.ErrPaymailDBFailed.Wrap(err)
	}

	pki, pkiDerivationKey, err := s.pki(ctx, paymailModel)
	if err != nil {
		return nil, err
	}

	dest, err := type42.NewDestinationWithRandomReference(pki)
	if err != nil {
		return nil, pmerrors.ErrPaymentDestination.Wrap(err)
	}

	address, err := script.NewAddressFromPublicKey(dest.PubKey, true)
	if err != nil {
		return nil, pmerrors.ErrPaymentDestination.Wrap(err)
	}

	lockingScript, err := p2pkh.Lock(address)
	if err != nil {
		return nil, pmerrors.ErrPaymentDestination.Wrap(err)
	}

	err = s.addresses.Create(ctx, &addressesmodels.NewAddress{
		UserID:  paymailModel.UserID,
		Address: address.AddressString,
		CustomInstructions: []bsv.CustomInstruction{
			{
				Type:        "type42",
				Instruction: pkiDerivationKey,
			},
			{
				Type:        "type42",
				Instruction: dest.DerivationKey,
			},
		},
	})
	if err != nil {
		return nil, pmerrors.ErrAddressSave.Wrap(err)
	}

	return &destinationData{
		address:       address.AddressString,
		lockingScript: lockingScript.String(),
		referenceID:   dest.ReferenceID,
	}, nil
}
