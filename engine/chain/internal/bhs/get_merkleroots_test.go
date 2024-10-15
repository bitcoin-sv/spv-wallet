package bhs_test

import (
	"context"
	"encoding/json"
	"net/url"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/engine/chain"
	"github.com/bitcoin-sv/spv-wallet/engine/chain/errors"
	"github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	"github.com/bitcoin-sv/spv-wallet/engine/tester"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/stretchr/testify/require"
)

func TestGetMerkleRootsSuccess(t *testing.T) {
	t.Run("Get MerkleRoots success", func(t *testing.T) {
		response := "{\"content\": [ {\"merkleRoot\": \"4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b\",\"blockHeight\": 0  },{  \"merkleRoot\": \"0e3e2357e806b6cdb1f70b54c3a3a17b6714ee1f0e68bebb44a74b1efd512098\", \"blockHeight\": 1}, {  \"merkleRoot\": \"9b0fc92260312ce44e74ef369f5c66bbb85848f2eddd5a7a1cde251e54ccfdd5\",   \"blockHeight\": 2  },  {   \"merkleRoot\": \"999e1c837c76a1b7fbb7e57baf87b309960f5ffefbf2a9b95dd890602272f644\", \"blockHeight\": 3 } ], \"page\": {   \"totalElements\": 866322,   \"size\": 4,   \"lastEvaluatedKey\": \"999e1c837c76a1b7fbb7e57baf87b309960f5ffefbf2a9b95dd890602272f644\" }}"
		httpClient := bhsMockMerkleRoots(200, response)

		var responseUnmarshalled *models.MerkleRootsBHSResponse
		err := json.Unmarshal([]byte(response), &responseUnmarshalled)
		require.NoError(t, err)

		service := chain.NewChainService(tester.Logger(t), httpClient, chainmodels.ARCConfig{}, bhsCfg(bhsURL, bhsToken))

		merkleroots, err := service.GetMerkleRoots(context.Background(), url.Values{})

		require.NoError(t, err)
		require.Equal(t, responseUnmarshalled, merkleroots)
	})
}

func TestGetMerkleRootsFailure(t *testing.T) {
	errorTestCases := map[string]struct {
		bhsToken         string
		bhsURL           string
		expectErr        error
		response         string
		responseCode     int
		batchSize        string
		lastEvaluatedKey string
	}{
		"Get MerkleRoots with wrong token": {
			bhsToken:         "wrong-token",
			bhsURL:           bhsURL,
			expectErr:        chainerrors.ErrBHSUnauthorized,
			responseCode:     400,
			response:         "{\"code\": \"error-unauthorized\",\"message\": \"unauthorized\"}",
			batchSize:        "",
			lastEvaluatedKey: "",
		},
		"Get MerkleRoots with wrong batch size": {
			bhsToken:         bhsToken,
			bhsURL:           bhsURL,
			expectErr:        chainerrors.ErrInvalidBatchSize,
			responseCode:     400,
			response:         "{\"code\": \"ErrInvalidBatchSize\",\"message\": \"batchSize must be 0 or a positive integer\"}",
			batchSize:        "-2",
			lastEvaluatedKey: "",
		},
		"Get MerkleRoots with invalid merkleroot": {
			bhsToken:         bhsToken,
			bhsURL:           bhsURL,
			expectErr:        chainerrors.ErrMerkleRootNotFound,
			responseCode:     404,
			response:         "{\"code\": \"ErrMerkleRootNotFound\",\"message\": \"No block with provided merkleroot was found\"}",
			batchSize:        "2",
			lastEvaluatedKey: "invalid-merkleroot",
		},
		"Get MerkleRoots with stale merkleroot": {
			bhsToken:         bhsToken,
			bhsURL:           bhsURL,
			expectErr:        chainerrors.ErrMerkleRootNotInLongestChain,
			responseCode:     409,
			response:         "{\"code\": \"ErrMerkleRootNotInLC\",\"message\": \"Provided merkleroot is not part of the longest chain\"}",
			batchSize:        "2",
			lastEvaluatedKey: "6ef51ab0d52991faf8f82e951a6ff9c40e6fdf1d56406067ff6882c6826323a5",
		},
		"Get MerkleRoots with bad url": {
			bhsToken:         bhsToken,
			bhsURL:           "%bad-url" + bhsURL,
			expectErr:        chainerrors.ErrBHSBadURL,
			responseCode:     0,
			response:         "",
			batchSize:        "",
			lastEvaluatedKey: "",
		},
	}

	for name, test := range errorTestCases {
		t.Run(name, func(t *testing.T) {
			httpClient := bhsMockMerkleRoots(test.responseCode, test.response)

			service := chain.NewChainService(tester.Logger(t), httpClient, chainmodels.ARCConfig{}, bhsCfg(test.bhsURL, test.bhsToken))
			urlQuery := url.Values{
				"lastEvaluatedKey": []string{test.lastEvaluatedKey},
				"batchSize":        []string{test.batchSize},
			}

			_, err := service.GetMerkleRoots(context.Background(), urlQuery)

			require.Error(t, err)
			require.ErrorIs(t, err, test.expectErr)
		})
	}
}
