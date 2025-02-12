package paymailmock

import (
	"net/http"

	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/jarcoal/httpmock"
)

// MockedP2PDestinationResponse is a mocked response for the P2P destinations endpoint
type MockedP2PDestinationResponse struct {
	Reference string
	Scripts   []string
	Satoshis  []bsv.Satoshis
}

var mockedLockingScripts = []string{
	"76a9143e2d1d795f8acaa7957045cc59376177eb04a3c588ac",
	"76a9145edbbedd1985b0e04423d40ac596e6104b0888a988ac",
	"76a9145afa2554bce1b8b83a87cadc7d6a654acf85954e88ac",
	"76a914f4afb9e287049c03d7588c2993084c95b44b92de88ac",
	"76a9145537ff4d8658ac8f8da4346e7e7abc64484a03d488ac",
}

// Responder returns a httpmock responder for the mocked P2P destinations response
func (m *MockedP2PDestinationResponse) Responder() httpmock.Responder {
	r, err := httpmock.NewJsonResponder(http.StatusOK, m.response())
	if err != nil {
		panic(spverrors.Wrapf(err, "cannot create mocked responder for P2P destinations"))
	}
	return r
}

func (m *MockedP2PDestinationResponse) response() obj {
	outs := make([]obj, len(m.Satoshis))
	for i, s := range m.Satoshis {
		outs[i] = obj{
			"script":   m.Scripts[i],
			"satoshis": s,
		}
	}
	return obj{
		"outputs":   outs,
		"reference": m.Reference,
	}
}

// P2PDestinationsForSats returns a mocked response for the P2P destinations endpoint
func P2PDestinationsForSats(satoshis bsv.Satoshis, moreSatoshis ...bsv.Satoshis) *MockedP2PDestinationResponse {
	outputSats := append([]bsv.Satoshis{satoshis}, moreSatoshis...)
	scripts := make([]string, len(outputSats))
	for i := range scripts {
		scripts[i] = mockedLockingScripts[i%len(mockedLockingScripts)]
	}

	return &MockedP2PDestinationResponse{
		Reference: "z0bac4ec-6f15-42de-9ef4-e60bfdabf4f7",
		Scripts:   scripts,
		Satoshis:  outputSats,
	}
}
