package testabilities

import (
	"encoding/json"
	"net/http"
	"slices"

	chainmodels "github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/bitcoin-sv/spv-wallet/models"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

type BlockHeadersServiceFixture interface {
	// WillRespondForMerkleRoots returns a http response for get merkleroots endpoint with
	// provided httpCode and response
	WillRespondForMerkleRoots(httpCode int, response string)

	WillRespondForMerkleRootsVerify(httpCode int, response *chainmodels.MerkleRootsConfirmations)
}

func (f *engineFixture) BHS() BlockHeadersServiceFixture {
	return f
}

func (f *engineFixture) WillRespondForMerkleRoots(httpCode int, response string) {
	responder := func(req *http.Request) (*http.Response, error) {
		res := httpmock.NewStringResponse(httpCode, response)
		res.Header.Set("Content-Type", "application/json")

		return res, nil
	}

	f.externalTransport.RegisterResponder("GET", "http://localhost:8080/api/v1/chain/merkleroot", responder)
}

func (f *engineFixture) WillRespondForMerkleRootsVerify(httpCode int, response *chainmodels.MerkleRootsConfirmations) {
	responder := func(req *http.Request) (*http.Response, error) {
		if response == nil {
			return httpmock.NewStringResponse(httpCode, ""), nil
		}
		return httpmock.NewJsonResponse(httpCode, response)
	}

	f.externalTransport.RegisterResponder("POST", "http://localhost:8080/api/v1/chain/merkleroot/verify", responder)
}

func (f *engineFixture) mockBHSGetMerkleRoots() {
	responder := func(req *http.Request) (*http.Response, error) {
		if req.Header.Get("Authorization") != "Bearer "+f.config.BHS.AuthToken {
			return httpmock.NewStringResponse(http.StatusUnauthorized, ""), nil
		}
		lastEvaluatedKey := req.URL.Query().Get("lastEvaluatedKey")
		merkleRootsRes, err := simulateBHSMerkleRootsAPI(lastEvaluatedKey)
		require.NoError(f.t, err)

		res := httpmock.NewStringResponse(http.StatusOK, merkleRootsRes)
		res.Header.Set("Content-Type", "application/json")

		return res, nil
	}

	f.externalTransport.RegisterResponder("GET", "http://localhost:8080/api/v1/chain/merkleroot", responder)
}

func simulateBHSMerkleRootsAPI(lastMerkleRoot string) (string, error) {
	var response models.MerkleRootsBHSResponse
	marshallResponseError := models.SPVError{StatusCode: http.StatusInternalServerError, Message: "Error during marshaling BHS response", Code: "err-marchall-bhs-res"}

	if lastMerkleRoot == "" {
		response.Content = fixtures.MockedBHSMerkleRootsData
		response.Page = models.ExclusiveStartKeyPageInfo{
			LastEvaluatedKey: "",
			TotalElements:    len(fixtures.MockedBHSMerkleRootsData),
			Size:             len(fixtures.MockedBHSMerkleRootsData),
		}

		resString, err := json.Marshal(response)
		if err != nil {
			return "", marshallResponseError.Wrap(err)
		}

		return string(resString), nil
	}

	lastMerkleRootIdx := slices.IndexFunc(fixtures.MockedBHSMerkleRootsData, func(mr models.MerkleRoot) bool {
		return mr.MerkleRoot == lastMerkleRoot
	})

	// handle case when lastMerkleRoot is already highest in the servers database
	if lastMerkleRootIdx == len(fixtures.MockedBHSMerkleRootsData)-1 {
		response.Content = []models.MerkleRoot{}
		response.Page = models.ExclusiveStartKeyPageInfo{
			LastEvaluatedKey: "",
			TotalElements:    len(fixtures.MockedBHSMerkleRootsData),
			Size:             0,
		}

		resString, err := json.Marshal(response)
		if err != nil {
			return "", marshallResponseError.Wrap(err)
		}

		return string(resString), nil
	}

	content := fixtures.MockedBHSMerkleRootsData[lastMerkleRootIdx+1:]
	lastEvaluatedKey := content[len(content)-1].MerkleRoot

	if lastEvaluatedKey == fixtures.MockedBHSMerkleRootsData[len(fixtures.MockedBHSMerkleRootsData)-1].MerkleRoot {
		lastEvaluatedKey = ""
	}

	response.Content = content
	response.Page = models.ExclusiveStartKeyPageInfo{
		LastEvaluatedKey: lastEvaluatedKey,
		TotalElements:    len(fixtures.MockedBHSMerkleRootsData),
		Size:             len(content),
	}

	resString, err := json.Marshal(response)
	if err != nil {
		return "", marshallResponseError.Wrap(err)
	}

	return string(resString), nil
}
