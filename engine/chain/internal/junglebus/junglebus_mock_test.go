package junglebus_test

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
)

const knownTx = "ea47e03186e59f8947d847e4eeaacde294a0a2db4d5e33b128430f2e2ee91015"
const wrongTxID = "wrong-txID"

func junglebusMockActivate(applyTimeout bool) *resty.Client {
	transport := httpmock.NewMockTransport()
	client := resty.New()
	client.GetClient().Transport = transport

	responder := func(status int, content string) func(req *http.Request) (*http.Response, error) {
		return func(req *http.Request) (*http.Response, error) {
			if applyTimeout {
				time.Sleep(100 * time.Millisecond)
			}
			res := httpmock.NewStringResponse(status, content)
			res.Header.Set("Content-Type", "application/json")
			return res, nil
		}
	}

	transport.RegisterResponder("GET", fmt.Sprintf("https://junglebus.gorillapool.io/v1/transaction/get/%s", knownTx), responder(http.StatusOK, `{
			"id": "ea47e03186e59f8947d847e4eeaacde294a0a2db4d5e33b128430f2e2ee91015",
			"transaction": "AQAAAALGRtoGYopYRrPBunmTa7CyugGezYrMIoFEG9GqRAjzzAAAAABqRzBEAiBQA1fX6FxiNAWv12fbEdLaILUMb8+XwNthxwhndI4+qAIgHF/vjX0fm08kpeXMuEBxqbH5gItdR1viHb5HXOFnAGBBIQIX/1jhAu02HUlGuvslevxySvPFCrEIpDM9WFroHyMAlf////9GMk6Kg/C83dcDWNjdjUJhNSncYzoUwAEjO8i+lZZC3wAAAABrSDBFAiEAr9QCPdR82p5/ZwT/+odA4poHpBoMXUiEBDR3u7YAgtkCICtJIese99j+goyCE/mOHL7Oosas3YZKUW7yl2nMZCO1QSEDIdOpsTxsCz0eFr3EwJkHScUss5T28G4xSxGGycU9YDr/////AhQAAAAAAAAAGXapFAWxMh7kR4+qCj0KroiTwUmBHsaOiKwTAAAAAAAAABl2qRQgrrV9SAnXS3u3jMu36wf8nFnhzIisAAAAAA==",
			"block_hash": "00000000000000000867325526c1f6d578dd535117f787fe2f67d78d1c5ccd7d",
			"block_height": 825476,
			"block_time": 1704283115,
			"block_index": 610,
			"addresses": [
				"1MyizGwAJJxxm35rN3yU3KTcW88DmqQ4de",
				"1NExg8ukZBtFoCTkNRsXaXKKo8owQRz2JM",
				"1X6ejopTnDoKxYxXgjnpMp8C7VuwfcyT1",
				"13yov8U51qN6WTZJ9EFnS2JvbRhhLsdDby"
			],
			"inputs": [],
			"outputs": [
				"76a91405b1321ee4478faa0a3d0aae8893c149811ec68e88ac",
				"76a91420aeb57d4809d74b7bb78ccbb7eb07fc9c59e1cc88ac"
			],
			"input_types": [],
			"output_types": [
				"pubkeyhash"
			],
			"contexts": [],
			"sub_contexts": [],
			"data": [],
			"merkle_proof": null
		}`),
	)

	transport.RegisterResponder(
		"GET",
		fmt.Sprintf("https://junglebus.gorillapool.io/v1/transaction/get/%s", wrongTxID),
		httpmock.NewStringResponder(http.StatusNotFound, `encoding/hex: odd length hex string`),
	)

	return client
}
