package paymail

import (
	"context"
	"fmt"

	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/go-paymail/server"
	"github.com/bitcoin-sv/go-paymail/spv"
	"github.com/bitcoin-sv/spv-wallet/engine/keys/type42"
	pmerrors "github.com/bitcoin-sv/spv-wallet/engine/paymail/errors"
)

// NewServiceProvider create a new paymail service server which handlers incoming paymail requests
func NewServiceProvider(repo Repository) server.PaymailServiceProvider {
	return &serviceProvider{
		repo: repo,
	}
}

type serviceProvider struct {
	repo Repository
}

func (s *serviceProvider) CreateAddressResolutionResponse(ctx context.Context, alias, domain string, senderValidation bool, metaData *server.RequestMetadata) (*paymail.ResolutionPayload, error) {
	//TODO implement me
	panic("implement me")
}

func (s *serviceProvider) CreateP2PDestinationResponse(ctx context.Context, alias, domain string, satoshis uint64, metaData *server.RequestMetadata) (*paymail.PaymentDestinationPayload, error) {
	//TODO implement me
	panic("implement me")
}

func (s *serviceProvider) GetPaymailByAlias(ctx context.Context, alias, domain string, metaData *server.RequestMetadata) (*paymail.AddressInformation, error) {
	model, err := s.repo.GetPaymailByAlias(alias, domain)
	if err != nil {
		return nil, pmerrors.ErrPaymailDBFailed.Wrap(err)
	}
	if model == nil {
		return nil, pmerrors.ErrPaymailNotFound
	}

	userPubKey, err := model.User.PubKeyObj()
	if err != nil {
		return nil, pmerrors.ErrPaymailPKI.Wrap(err)
	}

	pki, err := type42.PaymailPKI(userPubKey, model.Alias, model.Domain)
	if err != nil {
		return nil, pmerrors.ErrPaymailPKI.Wrap(err)
	}

	return &paymail.AddressInformation{
		Alias:  model.Alias,
		Avatar: model.AvatarURL,
		Domain: model.Domain,
		ID:     fmt.Sprintf("%d", model.ID),
		Name:   model.PublicName,
		PubKey: pki.ToDERHex(),
	}, nil
}

func (s *serviceProvider) RecordTransaction(ctx context.Context, p2pTx *paymail.P2PTransaction, metaData *server.RequestMetadata) (*paymail.P2PTransactionPayload, error) {
	//TODO implement me
	panic("implement me")
}

func (s *serviceProvider) VerifyMerkleRoots(ctx context.Context, merkleProofs []*spv.MerkleRootConfirmationRequestItem) error {
	//TODO implement me
	panic("implement me")
}
