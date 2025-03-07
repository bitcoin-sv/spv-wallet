package merkleroots_test

import (
	"encoding/json"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/stretchr/testify/require"
)

const merklerootsURL = "/api/v1/merkleroots"

type jsonObject = map[string]any

func TestGETMerkleRootsSuccess(t *testing.T) {
	testCases := map[string]struct {
		query            string
		expectedResponse jsonObject
	}{
		"Get MerkleRoots success no query params": {
			query: "",
			expectedResponse: jsonObject{
				"content": fixtures.MockedBHSMerkleRootsData,
				"page": jsonObject{
					"totalElements":    len(fixtures.MockedBHSMerkleRootsData),
					"size":             len(fixtures.MockedBHSMerkleRootsData),
					"lastEvaluatedKey": "",
				},
			},
		},
		"Get MerkleRoots success with last evaluated key param": {
			query: "?lastEvaluatedKey=df2b060fa2e5e9c8ed5eaf6a45c13753ec8c63282b2688322eba40cd98ea067a",
			expectedResponse: jsonObject{
				"content": fixtures.MockedBHSMerkleRootsData[5:],
				"page": jsonObject{
					"totalElements":    len(fixtures.MockedBHSMerkleRootsData),
					"size":             len(fixtures.MockedBHSMerkleRootsData[5:]),
					"lastEvaluatedKey": "",
				},
			},
		},
	}

	for name, tt := range testCases {
		t.Run(name, func(t *testing.T) {
			// given:
			expResponseJSON, err := json.Marshal(tt.expectedResponse)
			require.NoError(t, err, "Failed to marshall expected response")
			url := merklerootsURL
			given, then := testabilities.New(t)

			// and:
			cleanup := given.StartedSPVWallet()
			defer cleanup()
			client := given.HttpClient().ForUser()
			url = url + tt.query

			// when:
			res, _ := client.R().
				SetHeader("Content-Type", "application/json").
				Get(url)

			// then:
			then.Response(res).IsOK().WithJSONf(string(expResponseJSON))
		})
	}

}

func TestGETMerkleRootsFailure(t *testing.T) {
	testCases := map[string]struct {
		expectErr       string
		response        string
		responseCode    int
		expResponseCode int
	}{
		"Get MerkleRoots with wrong batch size": {
			responseCode: 400,
			response:     "{\"code\": \"ErrInvalidBatchSize\",\"message\": \"batchSize must be 0 or a positive integer\"}",
			expectErr:    "{\"code\":\"error-invalid-batch-size\",\"message\":\"batchSize must be 0 or a positive integer\"}",
		},
		"Get MerkleRoots with invalid merkleroot": {
			responseCode: 404,
			response:     "{\"code\": \"ErrMerkleRootNotFound\",\"message\": \"No block with provided merkleroot was found\"}",
			expectErr:    "{\"code\":\"error-merkleroot-not-found\",\"message\":\"No block with provided merkleroot was found\"}",
		},
		"Get MerkleRoots with stale merkleroot": {
			responseCode: 409,
			response:     "{\"code\": \"ErrMerkleRootNotInLC\",\"message\": \"Provided merkleroot is not part of the longest chain\"}",
			expectErr:    "{\"code\":\"error-merkleroot-not-part-of-longest-chain\",\"message\":\"Provided merkleroot is not part of the longest chain\"}",
		},
	}

	for name, tt := range testCases {
		t.Run(name, func(t *testing.T) {
			// given:
			given, then := testabilities.New(t)

			// and:
			cleanup := given.StartedSPVWallet()
			defer cleanup()
			client := given.HttpClient().ForUser()
			given.BHS().WillRespondForMerkleRoots(tt.responseCode, tt.response)

			// when:
			resErr := &models.ResponseError{}
			res, _ := client.R().
				SetHeader("Content-Type", "application/json").
				SetError(resErr).
				Get(merklerootsURL)

			// then:
			then.Response(res).HasStatus(tt.responseCode).WithJSONf(tt.expectErr)
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
