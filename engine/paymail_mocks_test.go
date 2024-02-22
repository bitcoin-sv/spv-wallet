package engine

import (
	"context"

	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/go-paymail/server"
	"github.com/bitcoin-sv/go-paymail/spv"
)

// Mock implementation of a service provider
type mockServiceProvider struct{}

// GetPaymailByAlias is a demo implementation of this interface
func (m *mockServiceProvider) GetPaymailByAlias(_ context.Context, _, _ string,
	_ *server.RequestMetadata,
) (*paymail.AddressInformation, error) {
	// Get the data from the demo database
	return nil, nil
}

// CreateAddressResolutionResponse is a demo implementation of this interface
func (m *mockServiceProvider) CreateAddressResolutionResponse(_ context.Context, _, _ string,
	_ bool, _ *server.RequestMetadata,
) (*paymail.ResolutionPayload, error) {
	// Generate a new destination / output for the basic address resolution
	return nil, nil
}

// CreateP2PDestinationResponse is a demo implementation of this interface
func (m *mockServiceProvider) CreateP2PDestinationResponse(_ context.Context, _, _ string,
	_ uint64, _ *server.RequestMetadata,
) (*paymail.PaymentDestinationPayload, error) {
	// Generate a new destination for the p2p request
	return nil, nil
}

// RecordTransaction is a demo implementation of this interface
func (m *mockServiceProvider) RecordTransaction(_ context.Context,
	_ *paymail.P2PTransaction, _ *server.RequestMetadata,
) (*paymail.P2PTransactionPayload, error) {
	// Record the tx into your datastore layer
	return nil, nil
}

// RecordTransaction is a demo implementation of this interface
func (m *mockServiceProvider) VerifyMerkleRoots(_ context.Context, _ []*spv.MerkleRootConfirmationRequestItem) error {
	// Verify merkle roots
	return nil
}
