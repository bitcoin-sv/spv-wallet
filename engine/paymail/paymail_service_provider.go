package paymail

import (
	"context"
	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/go-paymail/server"
	"github.com/bitcoin-sv/go-paymail/spv"
)

// NewServiceProvider create a new paymail service server which handlers incoming paymail requests
func NewServiceProvider() server.PaymailServiceProvider {
	return &serviceProvider{}
}

type serviceProvider struct {
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
	//TODO implement me
	panic("implement me")
}

func (s *serviceProvider) RecordTransaction(ctx context.Context, p2pTx *paymail.P2PTransaction, metaData *server.RequestMetadata) (*paymail.P2PTransactionPayload, error) {
	//TODO implement me
	panic("implement me")
}

func (s *serviceProvider) VerifyMerkleRoots(ctx context.Context, merkleProofs []*spv.MerkleRootConfirmationRequestItem) error {
	//TODO implement me
	panic("implement me")
}
