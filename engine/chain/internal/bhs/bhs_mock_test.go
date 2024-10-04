package bhs_test

import (
	"encoding/json"
	"fmt"
	"github.com/bitcoin-sv/go-paymail/spv"
	chainmodels "github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
	"net/http"
	"time"
)

const (
	bhsURL   = "http://localhost:8090"
	bhsToken = "mQZQ6WmxURxWz5ch"
)

func bhsMockVerify(response *chainmodels.MerkleRootsConfirmations, applyTimeout bool) *resty.Client {
	transport := httpmock.NewMockTransport()
	client := resty.New()
	client.GetClient().Transport = transport

	responder := func(req *http.Request) (*http.Response, error) {
		if applyTimeout {
			time.Sleep(100 * time.Millisecond)
		}
		if req.Header.Get("Authorization") != "Bearer "+bhsToken {
			return httpmock.NewStringResponse(http.StatusUnauthorized, ""), nil
		}
		var reqBody []*spv.MerkleRootConfirmationRequestItem
		_ = json.NewDecoder(req.Body).Decode(&reqBody)
		if len(reqBody) == 0 {
			return httpmock.NewStringResponse(http.StatusBadRequest, "at least one merkleroot is required"), nil
		}
		res, _ := httpmock.NewJsonResponse(http.StatusOK, response)
		res.Header.Set("Content-Type", "application/json")
		return res, nil
	}

	transport.RegisterResponder("POST", fmt.Sprintf("%s/api/v1/chain/merkleroot/verify", bhsURL), responder)

	return client
}

func bhsCfg(url, authToken string) chainmodels.BHSConfig {
	return chainmodels.BHSConfig{
		URL:       url,
		AuthToken: authToken,
	}
}
