package bhs_test

import (
	"context"
	"testing"
	"time"

	"github.com/bitcoin-sv/go-paymail/spv"
	"github.com/bitcoin-sv/spv-wallet/engine/chain"
	"github.com/bitcoin-sv/spv-wallet/engine/chain/errors"
	"github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	"github.com/bitcoin-sv/spv-wallet/engine/tester"
	"github.com/stretchr/testify/require"
)

/**
NOTE: switch httpClient to resty.New() tu call actual ARC server
*/

var validMerkleRootsReq = []*spv.MerkleRootConfirmationRequestItem{
	{
		MerkleRoot:  "4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b",
		BlockHeight: 0,
	},
	{
		MerkleRoot:  "0e3e2357e806b6cdb1f70b54c3a3a17b6714ee1f0e68bebb44a74b1efd512098",
		BlockHeight: 1,
	},
}

func TestVerifyMerkleRoots(t *testing.T) {
	tests := map[string]struct {
		request  []*spv.MerkleRootConfirmationRequestItem
		response *chainmodels.MerkleRootsConfirmations
	}{
		"Verify for CONFIRMED merkle roots": {
			request: validMerkleRootsReq,
			response: &chainmodels.MerkleRootsConfirmations{
				ConfirmationState: chainmodels.MRConfirmed,
				Confirmations: []chainmodels.MerkleRootConfirmation{
					{
						Hash:         "000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f",
						BlockHeight:  0,
						MerkleRoot:   "4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b",
						Confirmation: chainmodels.MRConfirmed,
					},
					{
						Hash:         "00000000839a8e6886ab5951d76f411475428afc90947ee320161bbf18eb6048",
						BlockHeight:  1,
						MerkleRoot:   "0e3e2357e806b6cdb1f70b54c3a3a17b6714ee1f0e68bebb44a74b1efd512098",
						Confirmation: chainmodels.MRConfirmed,
					},
				},
			},
		},
		"Verify with one wrong hash": {
			request: []*spv.MerkleRootConfirmationRequestItem{
				{
					MerkleRoot:  "wronge4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b",
					BlockHeight: 0,
				},
				{
					MerkleRoot:  "0e3e2357e806b6cdb1f70b54c3a3a17b6714ee1f0e68bebb44a74b1efd512098",
					BlockHeight: 1,
				},
			},
			response: &chainmodels.MerkleRootsConfirmations{
				ConfirmationState: chainmodels.MRInvalid,
				Confirmations: []chainmodels.MerkleRootConfirmation{
					{
						Hash:         "",
						BlockHeight:  0,
						MerkleRoot:   "wronge4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b",
						Confirmation: chainmodels.MRInvalid,
					},
					{
						Hash:         "00000000839a8e6886ab5951d76f411475428afc90947ee320161bbf18eb6048",
						BlockHeight:  1,
						MerkleRoot:   "0e3e2357e806b6cdb1f70b54c3a3a17b6714ee1f0e68bebb44a74b1efd512098",
						Confirmation: chainmodels.MRConfirmed,
					},
				},
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			httpClient := bhsMockVerify(test.response, false)

			service := chain.NewChainService(tester.Logger(t), httpClient, chainmodels.ARCConfig{}, bhsCfg(bhsURL, bhsToken))

			confirmations, err := service.VerifyMerkleRoots(context.Background(), test.request)

			require.NoError(t, err)
			require.NotNil(t, confirmations)
			require.Equal(t, test.response, confirmations)
		})
	}
}

func TestVerifyMerkleRootsErrorCases(t *testing.T) {
	errTestCases := map[string]struct {
		bhsToken  string
		bhsURL    string
		expectErr error
		request   []*spv.MerkleRootConfirmationRequestItem
	}{
		"Verify MR with wrong token": {
			bhsToken:  "wrong-token",
			bhsURL:    bhsURL,
			expectErr: chainerrors.ErrBHSUnauthorized,
			request:   validMerkleRootsReq,
		},
		"Verify MR endpoint wrong but reachable": {
			bhsToken:  bhsToken,
			bhsURL:    bhsURL + "/wrong",
			expectErr: chainerrors.ErrBHSUnreachable,
			request:   validMerkleRootsReq,
		},
		"Verify MR endpoint with wrong arcURL": {
			bhsToken:  bhsToken,
			bhsURL:    "wrong-url",
			expectErr: chainerrors.ErrBHSUnreachable,
			request:   validMerkleRootsReq,
		},
		"Verify MR endpoint with empty merkleroots": {
			bhsToken:  bhsToken,
			bhsURL:    bhsURL,
			expectErr: chainerrors.ErrBHSBadRequest,
			request:   []*spv.MerkleRootConfirmationRequestItem{},
		},
	}

	for name, test := range errTestCases {
		t.Run(name, func(t *testing.T) {
			httpClient := bhsMockVerify(&chainmodels.MerkleRootsConfirmations{}, false)

			service := chain.NewChainService(tester.Logger(t), httpClient, chainmodels.ARCConfig{}, bhsCfg(test.bhsURL, test.bhsToken))

			confirmations, err := service.VerifyMerkleRoots(context.Background(), test.request)

			require.Error(t, err)
			require.ErrorIs(t, err, test.expectErr)
			require.Nil(t, confirmations)
		})
	}
}

func TestVerifyMerkleRootsTimeouts(t *testing.T) {
	t.Run("VerifyMerkleRoots interrupted by ctx timeout", func(t *testing.T) {
		httpClient := bhsMockVerify(&chainmodels.MerkleRootsConfirmations{}, true)

		service := chain.NewChainService(tester.Logger(t), httpClient, chainmodels.ARCConfig{}, bhsCfg(bhsURL, bhsToken))

		ctx, cancel := context.WithTimeout(context.Background(), 1)
		defer cancel()

		confirmations, err := service.VerifyMerkleRoots(ctx, validMerkleRootsReq)

		require.Error(t, err)
		require.ErrorIs(t, err, chainerrors.ErrBHSUnreachable)
		require.ErrorIs(t, err, context.DeadlineExceeded)
		require.Nil(t, confirmations)
	})

	t.Run("VerifyMerkleRoots interrupted by resty timeout", func(t *testing.T) {
		httpClient := bhsMockVerify(&chainmodels.MerkleRootsConfirmations{}, true)
		httpClient.SetTimeout(1 * time.Millisecond)

		service := chain.NewChainService(tester.Logger(t), httpClient, chainmodels.ARCConfig{}, bhsCfg(bhsURL, bhsToken))

		confirmations, err := service.VerifyMerkleRoots(context.Background(), validMerkleRootsReq)

		require.Error(t, err)
		require.ErrorIs(t, err, chainerrors.ErrBHSUnreachable)
		require.ErrorIs(t, err, context.DeadlineExceeded)
		require.Nil(t, confirmations)
	})
}
