package junglebus_test

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
)

const knownTx1 = "ea47e03186e59f8947d847e4eeaacde294a0a2db4d5e33b128430f2e2ee91015"
const knownTx2 = "68fae29bb880e0c9fef78269a5199f46c1fc2fabe30d156825ac8db87aa70a48"
const unknownTx = "2010c88a4d8edc96296a27b35c5ca0ce2fe41d8d9af86040c14b3e2c55d44529"
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

	transport.RegisterResponder("GET", fmt.Sprintf("https://junglebus.gorillapool.io/v1/transaction/get/%s", knownTx1), responder(http.StatusOK, `{
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

	transport.RegisterResponder("GET", fmt.Sprintf("https://junglebus.gorillapool.io/v1/transaction/get/%s", knownTx2), responder(http.StatusOK, `{
			"id": "68fae29bb880e0c9fef78269a5199f46c1fc2fabe30d156825ac8db87aa70a48",
			"transaction": "AQAAAAJ3kIKpTsPBEzN6/uVme2nWTX11kUcskqTN432xVxnqRQEAAABqRzBEAiA47MKxSC3xF4cXIrY0okz4PRgAiia6LpjzQquLHfTf7QIgLgsemJcmE0kQATMpskTe6be45Io7outZJrs3XKZyM7JBIQLY0E7RwBkZ8Jc+EJHgVO23869GmSmuXhLQGvMvVZW1zP/////gce8B1l/L9rCILqIhameSiqBbb/o07NMvyVDpasrVIAEAAABqRzBEAiBdt7eaMHXnvru7rKqftYPhsAioLeqv1FEuqLhLsdTcowIgRKagoQJJlSEzO8uMpPNAJejAyA8715q75Y8aoK98cdtBIQI+zYsDDrI5kcxavnGKeOyAq7p0CmzRmQN7E/GrTgxzdv////8DFAAAAAAAAAAZdqkUl19u7zIpwIlq2i6tVQcCrz/aUSiIrAQQAAAAAAAAGXapFEhmM9LA4etvf6ecUlxOk9CotMIGiKxtAAAAAAAAABl2qRSQt/wq3+wdnaoGIvhXoG2bsS7iGIisAAAAAA==",
			"block_hash": "0000000000000000062cdbcb37156f75f3e7c7ceaf76fd9fbd286949491ee4d0",
			"block_height": 825612,
			"block_time": 1704370007,
			"block_index": 132,
			"addresses": [
				"1AJYBz1GvEMz5Ut56QbhwVrtGPP6dXpeCX",
				"1A1HscXc7Pqh1k5RaXQ33PBDn7cg8BJHXq",
				"1EoPMnseaoHoCeTHgX3rBtSf3mn8B4eN32",
				"17bp5xmGwLMG4E9sTpxJvAQGMnSe4fhC9g",
				"1ECCi8PyjBArzmq25vTWYSQn3y19GTgmod"
			],
			"inputs": [],
			"outputs": [
				"76a914975f6eef3229c0896ada2ead550702af3fda512888ac",
				"76a914486633d2c0e1eb6f7fa79c525c4e93d0a8b4c20688ac",
				"76a91490b7fc2adfec1d9daa0622f857a06d9bb12ee21888ac"
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
		httpmock.NewStringResponder(http.StatusNotFound, `"encoding/hex: invalid byte: U+0077 'w'"`),
	)

	transport.RegisterResponder(
		"GET",
		fmt.Sprintf("https://junglebus.gorillapool.io/v1/transaction/get/%s", unknownTx),
		httpmock.NewStringResponder(http.StatusNotFound, `"tx-not-found"`),
	)

	return client
}
