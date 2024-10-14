package arc_test

import (
	"fmt"
	"net/http"
	"time"

	"github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	"github.com/go-resty/resty/v2"
	"github.com/jarcoal/httpmock"
)

const (
	minedTxID         = "4dff1d32c1a02d7797e33d7c4ab2f96fe6699005b6d79e6391bdf5e358232e06"
	unknownTxID       = "aaaa1d32c1a02d7797e33d7c4ab2f96fe6699005b6d79e6391bdf5e358232e06"
	wrongButReachable = "/wrong/url"
	arcURL            = "https://arc.taal.com"
	arcToken          = "mainnet_06770f425eb00298839a24a49cbdc02c"
	invalidTxID       = "invalid"
)

// broadcast transaction cases
const (
	validRawHex         = "0100000001776f9d2c0d80b612ca54ebca1fa3bd38db375756ec7778edddc323569e06dc96010000006b483045022100f83415750880cc9464b752c33215ede7568c45c83cd3ccb841787edd1219d368022065e71b09241c889529e1b979c9cc9f569263386811f0005485b2617134722dd04121020d0ace627fbf80e20ce54d8bcfa5aa41f4a6f2d9c113bac0cf1b1316a96f0c9fffffffff0201000000000000001976a9147d71e2f8ce19409fb93312c9f019a7dee3d14ea188ac0c000000000000001976a914e5a4f06c1a24c2a9ddffcae4405aaf6c203ba3d288ac00000000"
	sourceOfValidRawHex = "0100000001f7af56f1d861954b4c2b5a55619e57d3ccf266a71ef9e3ff7e5adfefd8f85d30010000006a47304402200a9f2f7a1222329b97e23bc5c4def0823600495e2afa3bc6ae7be5d5bbd6aff20220525f69245f52548066c5319423578688ef46bab18b8185645dbe2587f0c9dbb341210237af69199d207ff5ce2309697644eac2babea83eadefd8d38384c8cae3e87cdfffffffff0201000000000000001976a9146fdeac955937ceec0493e6960421849636727fc088ac0e000000000000001976a9149c85e1657fb23d41e147124adef6df7933e86cb288ac00000000"
	efOfValidRawHex     = "010000000000000000ef01776f9d2c0d80b612ca54ebca1fa3bd38db375756ec7778edddc323569e06dc96010000006b483045022100f83415750880cc9464b752c33215ede7568c45c83cd3ccb841787edd1219d368022065e71b09241c889529e1b979c9cc9f569263386811f0005485b2617134722dd04121020d0ace627fbf80e20ce54d8bcfa5aa41f4a6f2d9c113bac0cf1b1316a96f0c9fffffffff0e000000000000001976a9149c85e1657fb23d41e147124adef6df7933e86cb288ac0201000000000000001976a9147d71e2f8ce19409fb93312c9f019a7dee3d14ea188ac0c000000000000001976a914e5a4f06c1a24c2a9ddffcae4405aaf6c203ba3d288ac00000000"

	// https://whatsonchain.com/tx/88a7c0ed1cb4767cfc8e7434561379eaea21ae78e480cacf4e69284387057c70
	txWithMultipleInputs                  = "01000000021b4ae503913172c5e16bd89dabb71d353c5b9cb2a1c69970fd4e690e49f97410010000006a47304402203127d53ed2ed8843d95ad0da49659e086e298dc8c4abf946656eeae1fd5c8c8602205e10f3bd2c3f01c08903c3969d138d62b1e76c96f856e2743cf3069cb4695a75412102792258b7fba50c8a1d6154f0b4be4a4e57b078efe1b47946c010697e99dde791ffffffffc88e4c870d61d7e14b8931d941c888ffb36ad58c52364e49c2df20f565dadecd010000006b483045022100c42531f0b50acab6fd1f63b30a2b1046ac29965bc0cf41409b180d3d4b91abec022006fbada6d4969de5f16297f8905a8a00763ca19a1a27ac31b1e77d128777a140412103d21e72986de0d354aff1dd737a066b6b786bc204bec22b3941e10e9575a7aa7bffffffff0214000000000000001976a914e8964298fcaa506f39e6d1d1f29657f79c1e72e788ac09000000000000001976a914e0bd3f2d5c1919109831bfad40b8eb293c07621b88ac00000000"
	sourceOneOfTxWithMultipleInputs       = "010000000124eebc416395164f0361f40aa2f555c26e9715d34e3c053d4e2a320465aebd71010000006a473044022018f346a2f9ef9b97d10b8771b5062a8dc689a7cf5c7dc75f4043bcf9e9c84aad022065dad45fc43b270ea6050c82a1dd28e2d03c5544e983a6b1c814a8bf76c0d79141210264250fb3346aaa01d758219d4c5707cafefe2224f4f78ee91eea50a054e5d704ffffffff0201000000000000001976a914e0842daa9d18a889c57d99aa510e5492c950bf9988ac10000000000000001976a914f8704d915ad7d2b559f61bad6c31b60deac52a3788ac00000000"
	txIDOfSourceTwoOfTxWithMultipleInputs = "cddeda65f520dfc2494e36528cd56ab3ff88c841d931894be1d7610d874c8ec8"
	efHexOfTxWithMultipleInputs           = "010000000000000000ef021b4ae503913172c5e16bd89dabb71d353c5b9cb2a1c69970fd4e690e49f97410010000006a47304402203127d53ed2ed8843d95ad0da49659e086e298dc8c4abf946656eeae1fd5c8c8602205e10f3bd2c3f01c08903c3969d138d62b1e76c96f856e2743cf3069cb4695a75412102792258b7fba50c8a1d6154f0b4be4a4e57b078efe1b47946c010697e99dde791ffffffff10000000000000001976a914f8704d915ad7d2b559f61bad6c31b60deac52a3788acc88e4c870d61d7e14b8931d941c888ffb36ad58c52364e49c2df20f565dadecd010000006b483045022100c42531f0b50acab6fd1f63b30a2b1046ac29965bc0cf41409b180d3d4b91abec022006fbada6d4969de5f16297f8905a8a00763ca19a1a27ac31b1e77d128777a140412103d21e72986de0d354aff1dd737a066b6b786bc204bec22b3941e10e9575a7aa7bffffffff0e000000000000001976a9146b8297b1c3cd9ec13151c90d29e3a96f147535a688ac0214000000000000001976a914e8964298fcaa506f39e6d1d1f29657f79c1e72e788ac09000000000000001976a914e0bd3f2d5c1919109831bfad40b8eb293c07621b88ac00000000"

	fallbackRawHex = "010000000116d60a1563239eac2295b4eecbc6982ff6d007f480e52505c78f803bc8e03a05010000006a473044022024f84674219f2ec2fb78d38bcd19d4ae5b44dd45474d7680d56662a56b127326022025590d4aec95942b0eb6d52e679e4c98939d7a72b5901fae46354552af42cdeb412103ec9a56e27b5b773459c7cef92683a0498da7073346728a724d1878a9d7ce9615ffffffff0201000000000000001976a9149eb8198a2f08551afc193663a0dd80a9ed2f3c1288ac10000000000000001976a914098d21f508a39588d31dd746757c83b7d790cccc88ac00000000"

	oldWithDoubleSpentHex = "0100000001293f17ea61f50d5ea815780c3d571f0f475533b8e812189724ab8e14b77e1616000000006a4730440220607cf28232c23fc7a2283e89466e740a02dcdc6bc5fe094a281ac89d81c1987f02206baadd4c704c0b8099e1df85ca90cce74023fef6f6381ca93f422f9ac5af4d58412103513000984c44b7316671c1875c32eaeeacfd886f561623479794913c1cb91f73ffffffff01000000000000000038006a35323032342d31302d31342030373a34323a30362e31343333383534202b303230302043455354206d3d2b302e30313833333730303100000000"
	newWithDoubleSpentHex = "0100000001ef9af77a38ee871bcca33df1260ab0b5c647743b4da33e417c4986150af6131b000000006a47304402207e93d325c42536b8255db76be8e34dfb486563f24e251a38f4afd3172bb78295022003ff04493bda06d020237f560f1ea0625da111a2c1bed72d0d1ae602dbaa12984121024a1ffa7ae3125b10c870d883d9dcd256f7b5ac51902f6901138e0861f95a9f59ffffffff01000000000000000038006a35323032342d31302d31342030383a30383a35382e37323636363431202b303230302043455354206d3d2b302e30313839303039303100000000"
	malformedTxHex        = "0100000001ef9af77a38ee871bcca33df1260ab0b5c647743b4da33e417c4986150af6131b0000000000ffffffff01000000000000000038006a35323032342d31302d31342030383a31313a33342e37313032333437202b303230302043455354206d3d2b302e30313832363439303100000000"
)

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

	transport.RegisterResponder("GET", fmt.Sprintf("%s/v1/policy", arcURL), responder(http.StatusOK, `{
			"policy": {
				"maxscriptsizepolicy": 100000000,
				"maxtxsigopscountspolicy": 4294967295,
				"maxtxsizepolicy": 100000000,
				"miningFee": {
					"bytes": 1000,
					"satoshis": 1
				}
			},
			"timestamp": "2024-10-02T07:36:33.589144918Z"
		}`),
	)

	transport.RegisterResponder("GET", arcURL+wrongButReachable, responder(http.StatusNotFound, `{
			"message": "no matching operation was found"
		}`),
	)

	transport.RegisterMatcherResponder("POST", fmt.Sprintf("%s/v1/tx", arcURL),
		httpmock.BodyContainsString(efOfValidRawHex),
		responder(http.StatusOK, `{
			"blockHash": "",
			"blockHeight": 0,
			"competingTxs": null,
			"extraInfo": "",
			"merklePath": "",
			"timestamp": "2024-09-27T06:11:41.417057192Z",
			"txStatus": "SEEN_ON_NETWORK",
			"txid": "2978f03c8a21bf90b5980113f988c39ef4ae691b9bedd5178c50ebb9c034dabf"
		}`),
	)

	transport.RegisterMatcherResponder("POST", fmt.Sprintf("%s/v1/tx", arcURL),
		httpmock.BodyContainsString(efHexOfTxWithMultipleInputs),
		responder(http.StatusOK, `{
			"blockHash": "",
			"blockHeight": 0,
			"competingTxs": null,
			"extraInfo": "",
			"merklePath": "",
			"timestamp": "2024-09-27T06:11:41.417057192Z",
			"txStatus": "SEEN_ON_NETWORK",
			"txid": "88a7c0ed1cb4767cfc8e7434561379eaea21ae78e480cacf4e69284387057c70"
		}`),
	)

	transport.RegisterResponder("GET", fmt.Sprintf("https://junglebus.gorillapool.io/v1/transaction/get/%s", txIDOfSourceTwoOfTxWithMultipleInputs), httpmock.NewStringResponder(http.StatusOK, `{
			"id": "cddeda65f520dfc2494e36528cd56ab3ff88c841d931894be1d7610d874c8ec8",
			"transaction": "AQAAAAHU58f2jMJt3XzGjJEKINLPVzwd2Mr6NDEAq8exla/vIgEAAABrSDBFAiEA3rvUh3L5fGG8nzMdxTW6AoKarzlehm3pHMDDULQ+f0sCIAmo1o/v9WUJD62kTZgsZ3iBYn3AjpkjOG7iWyedxxCxQSEDXI/Xt/qQrisBpMkdoNh/87u8M5DZ3me2n61SqLeP9J3/////AgEAAAAAAAAAGXapFAS8COAvcQwoaykycYzP1nGgyBZEiKwOAAAAAAAAABl2qRRrgpexw82ewTFRyQ0p46lvFHU1poisAAAAAA=="
		}`).HeaderSet(http.Header{"Content-Type": []string{"application/json"}}),
	)

	transport.RegisterMatcherResponder("POST", fmt.Sprintf("%s/v1/tx", arcURL),
		httpmock.BodyContainsString(fallbackRawHex),
		responder(http.StatusOK, `{
			"blockHash": "",
			"blockHeight": 0,
			"competingTxs": null,
			"extraInfo": "",
			"merklePath": "",
			"timestamp": "2024-09-27T06:11:41.417057192Z",
			"txStatus": "SEEN_ON_NETWORK",
			"txid": "305df8d8efdf5a7effe3f91ea766f2ccd3579e61555a2b4c4b9561d8f156aff7"
		}`),
	)

	transport.RegisterMatcherResponder("POST", fmt.Sprintf("%s/v1/tx", arcURL),
		httpmock.BodyContainsString(oldWithDoubleSpentHex),
		responder(http.StatusOK, `{
			"blockHash" : "",
			"blockHeight" : 0,
			"competingTxs" : null,
			"extraInfo" : "",
			"merklePath" : "",
			"status" : 200,
			"timestamp" : "2024-10-14T06:03:29.085353Z",
			"title" : "OK",
			"txStatus" : "SEEN_IN_ORPHAN_MEMPOOL",
			"txid" : "65e965d0ff776bf6dbbcc257d62a9cd7b52bd4caee5999e65fc83656550e2756"
		}`),
	)

	transport.RegisterMatcherResponder("POST", fmt.Sprintf("%s/v1/tx", arcURL),
		httpmock.BodyContainsString(newWithDoubleSpentHex),
		responder(http.StatusOK, `{
		  "blockHash": "",
		  "blockHeight": 0,
		  "competingTxs": [
			"62bf0fad6d45a7fdfbb2aae58a99c2b0812b1fa7141c4f98087ad721e3590731"
		  ],
		  "extraInfo": "",
		  "merklePath": "",
		  "status": 200,
		  "timestamp": "2024-10-11T11:06:12.891372826Z",
		  "title": "OK",
		  "txStatus": "DOUBLE_SPEND_ATTEMPTED",
		  "txid": "4997c2b412a9b9ae82074ef41f561371c74a33ff01cefce75b56caf546a77d19"
		}`),
	)

	transport.RegisterMatcherResponder("POST", fmt.Sprintf("%s/v1/tx", arcURL),
		httpmock.BodyContainsString(malformedTxHex),
		responder(461, `{
			"detail": "Transaction is malformed and cannot be processed",
			"extraInfo": "arc error 461: script execution failed\nindex 0 is invalid for stack size 0",
			"instance": null,
			"status": 461,
			"title": "Malformed transaction",
			"txid": "5b5ead5a42c4320f5b345e387c3648e8e789b171d17002e3efdc828534979f57",
			"type": "https://bitcoin-sv.github.io/arc/#/errors?id=_461"
		}`),
	)

	return client
}

func arcCfg(url, token string) chainmodels.ARCConfig {
	return chainmodels.ARCConfig{
		URL:          url,
		Token:        token,
		DeploymentID: "spv-wallet-test-arc-connection",
	}
}
