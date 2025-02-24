package paymailmock

import (
	"encoding/json"
	"net/http"

	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/jarcoal/httpmock"
	"github.com/samber/lo"
)

// MockedP2PDestination is a mocked response for the P2P destinations endpoint
type MockedP2PDestination struct {
	CustomDistribution []bsv.Satoshis
}

var mockedLockingScripts = []string{
	"76a9143e2d1d795f8acaa7957045cc59376177eb04a3c588ac",
	"76a9145edbbedd1985b0e04423d40ac596e6104b0888a988ac",
	"76a9145afa2554bce1b8b83a87cadc7d6a654acf85954e88ac",
	"76a914f4afb9e287049c03d7588c2993084c95b44b92de88ac",
	"76a9145537ff4d8658ac8f8da4346e7e7abc64484a03d488ac",
}

// MockedReferenceID is returned always by the mocked P2P destinations response
const MockedReferenceID = "z0bac4ec-6f15-42de-9ef4-e60bfdabf4f7"

// MockedP2PDestinationOutput is a part of the mocked P2P destinations response
type MockedP2PDestinationOutput struct {
	Script   string       `json:"script"`
	Satoshis bsv.Satoshis `json:"satoshis"`
}

// MockedP2PDestinationResponse is a model for the mocked P2P destinations response
type MockedP2PDestinationResponse struct {
	Reference string                       `json:"reference"`
	Outputs   []MockedP2PDestinationOutput `json:"outputs"`
}

// Responder returns a httpmock responder for the mocked P2P destinations response
func (m *MockedP2PDestination) Responder() httpmock.Responder {
	return func(request *http.Request) (*http.Response, error) {
		var payload paymail.PaymentRequest
		err := json.NewDecoder(request.Body).Decode(&payload)
		if err != nil {
			return httpmock.NewStringResponse(http.StatusBadRequest, "invalid json"), nil
		}

		var distribution []bsv.Satoshis
		if len(m.CustomDistribution) == 0 {
			distribution = []bsv.Satoshis{bsv.Satoshis(payload.Satoshis)}
		} else {
			distribution = m.CustomDistribution
		}

		outputs := lo.Map(distribution, func(sats bsv.Satoshis, i int) MockedP2PDestinationOutput {
			return MockedP2PDestinationOutput{
				Script:   mockedLockingScripts[i%len(mockedLockingScripts)],
				Satoshis: sats,
			}
		})

		r, err := httpmock.NewJsonResponse(http.StatusOK, MockedP2PDestinationResponse{
			Outputs:   outputs,
			Reference: MockedReferenceID,
		})
		if err != nil {
			panic(spverrors.Wrapf(err, "cannot create mocked responder for record tx response"))
		}

		return r, nil
	}
}

// P2PDestinationResponse returns a new mocked response for the P2P destinations endpoint
func P2PDestinationResponse() *MockedP2PDestination {
	return &MockedP2PDestination{}
}

// P2PDestinationsForSats returns a mocked response for the P2P destinations endpoint
func P2PDestinationsForSats(satoshis bsv.Satoshis, moreSatoshis ...bsv.Satoshis) *MockedP2PDestination {
	return &MockedP2PDestination{
		CustomDistribution: append([]bsv.Satoshis{satoshis}, moreSatoshis...),
	}
}
