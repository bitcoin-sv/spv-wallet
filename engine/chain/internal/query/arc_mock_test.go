package query_test

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
)

const minedTxID = "4dff1d32c1a02d7797e33d7c4ab2f96fe6699005b6d79e6391bdf5e358232e06"
const unknownTxID = "aaaa1d32c1a02d7797e33d7c4ab2f96fe6699005b6d79e6391bdf5e358232e06"
const wrongButReachable = "/wrong/url"
const arcURL = "https://arc.taal.com"
const arcToken = "mainnet_06770f425eb00298839a24a49cbdc02c"
const invalidTxID = "invalid"

func arcMockActivate(applyTimeout bool) *resty.Client {
	transport := httpmock.NewMockTransport()
	client := resty.New()
	client.GetClient().Transport = transport

	responder := func(status int, content string) func(req *http.Request) (*http.Response, error) {
		return func(req *http.Request) (*http.Response, error) {
			if applyTimeout {
				time.Sleep(100 * time.Millisecond)
			}
			if req.Header.Get("Authorization") != arcToken {
				return httpmock.NewStringResponse(http.StatusUnauthorized, ""), nil
			}
			res := httpmock.NewStringResponse(status, content)
			res.Header.Set("Content-Type", "application/json")
			return res, nil
		}
	}

	transport.RegisterResponder("GET", fmt.Sprintf("%s/v1/tx/%s", arcURL, minedTxID), responder(http.StatusOK, `{
			"blockHash": "0000000000000000034df47d8fe84ccf10267b4f6bc43be513d4604229d1c209",
			"blockHeight": 862510,
			"competingTxs": null,
			"extraInfo": "",
			"merklePath": "fe2e290d00080231006449ce1869e63013f9b3ad17151fe0fe37091c47fd9a70e03dddeb6a64a5592c3002062e2358e3f5bd91639ed7b6059069e66ff9b24a7c3de397772da0c1321dff4d011900bcab35ce0c50582723db10783f8b48f1f3165203ffb0644b91fd0d6cb4d6190f010d003e1c17e035a1248377cf3371863c853892283ef032abc53d8427b0b196368aec010700cf74d836ab526d6b3ef705f4adc1121724b886f6c6a79ccf080c1a0bbce712570102005a7c6ac761a529dc616656cf187b354516752372823f574706c61747741ac3d4010000d58c14f52200fd9d9e2ef87c993c99c2f28a636ebbfe88a0097066b5f10bc3a5010100e74bf2106c1b378d72b2e8e4f82646f955e0b6b9955505f7f3cddebce3ab733801010056365352ba7e5578ff8249905d25e272c540472276163b27b0e9c6d4d26b7d0e",
			"timestamp": "2024-09-27T06:11:41.417057192Z",
			"txStatus": "MINED",
			"txid": "4dff1d32c1a02d7797e33d7c4ab2f96fe6699005b6d79e6391bdf5e358232e06"
		}`),
	)

	transport.RegisterResponder("GET", fmt.Sprintf("%s/v1/tx/%s", arcURL, unknownTxID), responder(http.StatusNotFound, `{
			"detail": "The requested resource could not be found",
			"extraInfo": "transaction not found",
			"instance": null,
			"status": 404,
			"title": "Not found",
			"txid": null,
			"type": "https://bitcoin-sv.github.io/arc/#/errors?id=_404"
		}`),
	)

	transport.RegisterResponder("GET", arcURL+wrongButReachable, responder(http.StatusNotFound, `{
			"message": "no matching operation was found"
		}`),
	)

	transport.RegisterResponder("GET", fmt.Sprintf("%s/v1/tx/%s", arcURL, invalidTxID), responder(http.StatusConflict, `{
			"detail": "Transaction could not be processed",
			"extraInfo": "rpc error: code = Unknown desc = encoding/hex: invalid byte: U+0073 's'",
			"instance": null,
			"status": 409,
			"title": "Generic error",
			"txid": null,
			"type": "https://bitcoin-sv.github.io/arc/#/errors?id=_409"
		}`),
	)

	return client
}
