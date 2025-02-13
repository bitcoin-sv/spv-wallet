package paymailmock

import (
	"encoding/json"
	"github.com/bitcoin-sv/go-paymail"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/samber/lo"
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/jarcoal/httpmock"
)

// MockedP2PDestinationResponse is a mocked response for the P2P destinations endpoint
type MockedP2PDestinationResponse struct {
	CustomDistribution []bsv.Satoshis
	ExtraOutputs       int
}

var mockedLockingScripts = []string{
	"76a9143e2d1d795f8acaa7957045cc59376177eb04a3c588ac",
	"76a9145edbbedd1985b0e04423d40ac596e6104b0888a988ac",
	"76a9145afa2554bce1b8b83a87cadc7d6a654acf85954e88ac",
	"76a914f4afb9e287049c03d7588c2993084c95b44b92de88ac",
	"76a9145537ff4d8658ac8f8da4346e7e7abc64484a03d488ac",
}

const ReferenceID = "z0bac4ec-6f15-42de-9ef4-e60bfdabf4f7"

// Responder returns a httpmock responder for the mocked P2P destinations response
func (m *MockedP2PDestinationResponse) Responder() httpmock.Responder {
	return func(request *http.Request) (*http.Response, error) {
		var payload paymail.PaymentRequest
		err := json.NewDecoder(request.Body).Decode(&payload)
		if err != nil {
			return httpmock.NewStringResponse(http.StatusBadRequest, "invalid json"), nil
		}

		var distribution []bsv.Satoshis
		if len(m.CustomDistribution) == 0 {
			distribution = distribute(1+m.ExtraOutputs, bsv.Satoshis(payload.Satoshis))
		} else {
			distribution = m.CustomDistribution
		}

		outputs := lo.Map(distribution, func(sats bsv.Satoshis, i int) map[string]any {
			return map[string]any{
				"script":   mockedLockingScripts[i%len(mockedLockingScripts)],
				"satoshis": sats,
			}
		})

		r, err := httpmock.NewJsonResponse(http.StatusOK, map[string]any{
			"outputs":   outputs,
			"reference": ReferenceID,
		})
		if err != nil {
			panic(spverrors.Wrapf(err, "cannot create mocked responder for record tx response"))
		}

		return r, nil
	}
}

// P2PDestinationResponse returns a new mocked response for the P2P destinations endpoint
func P2PDestinationResponse() *MockedP2PDestinationResponse {
	return &MockedP2PDestinationResponse{}
}

// P2PDestinationsForSats returns a mocked response for the P2P destinations endpoint
func P2PDestinationsForSats(satoshis bsv.Satoshis, moreSatoshis ...bsv.Satoshis) *MockedP2PDestinationResponse {
	return &MockedP2PDestinationResponse{
		CustomDistribution: append([]bsv.Satoshis{satoshis}, moreSatoshis...),
	}
}

// P2PDestinationWithExtraOutputs returns a mocked response for the P2P destinations endpoint
func P2PDestinationWithExtraOutputs(outputs int) *MockedP2PDestinationResponse {
	return &MockedP2PDestinationResponse{
		ExtraOutputs: outputs,
	}
}

func distribute(items int, value bsv.Satoshis) []bsv.Satoshis {
	result := make([]bsv.Satoshis, items)
	if items <= 0 {
		return result
	}
	// If we have at least one unit per item,
	// use integer division and put any remainder in the first item.
	if value >= bsv.Satoshis(items) {
		base := value / bsv.Satoshis(items)
		rem := value % bsv.Satoshis(items)
		result[0] = base + rem
		for i := 1; i < items; i++ {
			result[i] = base
		}
	} else {
		// When value is less than items,
		// we can only assign a 1 to as many items as possible.
		// (distributed values cannot be less than 1 if nonzero)
		for i := 0; i < int(value); i++ {
			result[i] = 1
		}
		// The rest remain zero.
	}
	return result
}
