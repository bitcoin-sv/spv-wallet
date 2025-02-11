package testabilities

import (
	"net/http"

	chainmodels "github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

type ARCFixture interface {
	// WillRespondForBroadcast returns a http response for a broadcast request.
	WillRespondForBroadcast(httpCode int, info *chainmodels.TXInfo)

	// WillRespondForBroadcastWithSeenOnNetwork is a default ARC behaviour for broadcasting (happy path).
	WillRespondForBroadcastWithSeenOnNetwork(txID string)
}

func (f *engineFixture) ARC() ARCFixture {
	return f
}

func (f *engineFixture) WillRespondForBroadcast(httpCode int, info *chainmodels.TXInfo) {
	responder := func(req *http.Request) (*http.Response, error) {
		res, err := httpmock.NewJsonResponse(httpCode, info)
		require.NoError(f.t, err)
		res.Header.Set("Content-Type", "application/json")

		return res, nil
	}

	f.externalTransport.RegisterResponder("POST", "https://arc.taal.com/v1/tx", responder)
}

func (f *engineFixture) WillRespondForBroadcastWithSeenOnNetwork(txID string) {
	f.WillRespondForBroadcast(http.StatusOK, &chainmodels.TXInfo{
		TxID:     txID,
		TXStatus: chainmodels.SeenOnNetwork,
	})
}
