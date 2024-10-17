package merkleroots_test

import (
	"encoding/json"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/chain/errors"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/stretchr/testify/require"
)

const merklerootsURL = "/api/v1/merkleroots"

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
			expResponseJSON, _ := json.Marshal(tt.expectedResponse)
			given, then := testabilities.New(t)
			cleanup := given.StartedSPVWallet()
			defer cleanup()
			url := merklerootsURL

			if tt.query != "" {
				url = url + tt.query
			}

			// and
			client := given.HttpClient().ForUser()

			// when
			res, err := client.R().
				SetHeader("Content-Type", "application/json").
				Get(url)

			// then
			then.Response(res).IsOK().WithJSONf(string(expResponseJSON))
			require.NoError(t, err)
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
			given, then := testabilities.New(t)
			cleanup := given.StartedSPVWallet()
			defer cleanup()

			// and
			given.BHS().WillRespondForMerkleRoots(tt.responseCode, tt.response)
			client := given.HttpClient().ForUser()

			// when
			resErr := &models.ResponseError{}
			res, _ := client.R().
				SetHeader("Content-Type", "application/json").
				SetError(resErr).
				Get(tt.bhsURL)

			// then
			then.Response(res).IsNotSuccess()
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
		res, _ := client.R().Get(merklerootsURL)

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
		res, _ := client.R().Get(merklerootsURL)

		// then:
		then.Response(res).IsUnauthorizedForAdmin()
	})
}
