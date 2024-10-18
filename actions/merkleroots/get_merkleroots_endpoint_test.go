package merkleroots_test

import (
	"encoding/json"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/chain/errors"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/go-resty/resty/v2"
	"github.com/stretchr/testify/require"
)

const merklerootsURL = "/api/v1/merkleroots"

func setupWalletAndClientForUser(t *testing.T, responseCode *int, response *string) (*testabilities.SPVWalletApplicationAssertions, func(), *resty.Client) {
	given, then := testabilities.New(t)
	cleanup := given.StartedSPVWallet()
	client := given.HttpClient().ForUser()

	if response != nil && responseCode != nil {
		given.BHS().WillRespondForMerkleRoots(*responseCode, *response)
	}

	return &then, cleanup, client
}

func TestGETMerkleRootsSuccess(t *testing.T) {

	testCases := map[string]struct {
		query            string
		expectedResponse models.MerkleRootsBHSResponse
	}{
		"Get MerkleRoots success no query params": {
			query: "",
			expectedResponse: models.MerkleRootsBHSResponse{
				Content: testabilities.MockedBHSData,
				Page: models.ExclusiveStartKeyPageInfo{
					TotalElements:    len(testabilities.MockedBHSData),
					Size:             len(testabilities.MockedBHSData),
					LastEvaluatedKey: "",
				},
			},
		},
		"Get MerkleRoots success with last evaluated key param": {
			query: "?lastEvaluatedKey=df2b060fa2e5e9c8ed5eaf6a45c13753ec8c63282b2688322eba40cd98ea067a",
			expectedResponse: models.MerkleRootsBHSResponse{
				Content: testabilities.MockedBHSData[5:],
				Page: models.ExclusiveStartKeyPageInfo{
					TotalElements:    len(testabilities.MockedBHSData),
					Size:             len(testabilities.MockedBHSData[5:]),
					LastEvaluatedKey: "",
				},
			},
		},
	}

	for name, tt := range testCases {
		t.Run(name, func(t *testing.T) {

			// given
			expResponseJSON, err := json.Marshal(tt.expectedResponse)
			require.NoError(t, err, "Failed to marshall expected response")
			then, cleanup, client := setupWalletAndClientForUser(t, nil, nil)
			defer cleanup()
			url := merklerootsURL

			if tt.query != "" {
				url = url + tt.query
			}

			// when
			res, err := client.R().
				SetHeader("Content-Type", "application/json").
				Get(url)
			require.NoError(t, err, "Unexpected error occured while getting MerkleRoots")

			// then
			(*then).Response(res).IsOK().WithJSONf(string(expResponseJSON))
		})
	}

}

func TestGETMerkleRootsFailure(t *testing.T) {

	testCases := map[string]struct {
		bhsToken         string
		bhsURL           string
		expectErr        models.SPVError
		response         string
		responseCode     int
		batchSize        string
		lastEvaluatedKey string
	}{
		"Get MerkleRoots with wrong batch size": {
			bhsURL:       merklerootsURL,
			expectErr:    chainerrors.ErrInvalidBatchSize,
			responseCode: 400,
			response:     "{\"code\": \"ErrInvalidBatchSize\",\"message\": \"batchSize must be 0 or a positive integer\"}",
		},
		"Get MerkleRoots with invalid merkleroot": {
			bhsURL:       merklerootsURL,
			expectErr:    chainerrors.ErrMerkleRootNotFound,
			responseCode: 404,
			response:     "{\"code\": \"ErrMerkleRootNotFound\",\"message\": \"No block with provided merkleroot was found\"}",
		},
		"Get MerkleRoots with stale merkleroot": {
			bhsURL:       merklerootsURL,
			expectErr:    chainerrors.ErrMerkleRootNotInLongestChain,
			responseCode: 409,
			response:     "{\"code\": \"ErrMerkleRootNotInLC\",\"message\": \"Provided merkleroot is not part of the longest chain\"}",
		},
	}

	for name, tt := range testCases {
		t.Run(name, func(t *testing.T) {
			// given
			then, cleanup, client := setupWalletAndClientForUser(t, &tt.responseCode, &tt.response)
			defer cleanup()

			// when
			resErr := &models.ResponseError{}
			res, err := client.R().
				SetHeader("Content-Type", "application/json").
				SetError(resErr).
				Get(tt.bhsURL)
			require.NoError(t, err, "Unexpected error occured while getting MerkleRoots")

			// then
			(*then).Response(res).IsNotSuccess()
			require.NotNil(t, resErr)
			require.Equal(t, tt.expectErr.GetCode(), resErr.Code)
			require.Equal(t, tt.expectErr.GetMessage(), resErr.Message)
		})
	}

	t.Run("not allowed for anonymous", func(t *testing.T) {
		// given:
		given, then := testabilities.New(t)
		cleanup := given.StartedSPVWallet()
		defer cleanup()

		// and:
		client := given.HttpClient().ForAnonymous()

		// when:
		res, err := client.R().Get(merklerootsURL)
		require.NoError(t, err, "Unexpected error occured while getting MerkleRoots")

		// then:
		then.Response(res).IsUnauthorized()
	})

	t.Run("not allowed for admin", func(t *testing.T) {
		// given:
		given, then := testabilities.New(t)
		cleanup := given.StartedSPVWallet()
		defer cleanup()

		// and:
		client := given.HttpClient().ForAdmin()

		// when:
		res, err := client.R().Get(merklerootsURL)
		require.NoError(t, err, "Unexpected error occured while getting MerkleRoots")

		// then:
		then.Response(res).IsUnauthorizedForAdmin()
	})
}
