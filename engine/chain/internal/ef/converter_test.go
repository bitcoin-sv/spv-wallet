package ef_test

import (
	"context"
	"encoding/base64"
	sdk "github.com/bitcoin-sv/go-sdk/transaction"
	"github.com/bitcoin-sv/spv-wallet/engine/chain/internal/ef"
	"github.com/stretchr/testify/require"
	"testing"
)

const (
	// https://whatsonchain.com/tx/2978f03c8a21bf90b5980113f988c39ef4ae691b9bedd5178c50ebb9c034dabf
	tx1       = "AQAAAAF3b50sDYC2EspU68ofo7042zdXVux3eO3dwyNWngbclgEAAABrSDBFAiEA+DQVdQiAzJRkt1LDMhXt51aMRcg808y4QXh+3RIZ02gCIGXnGwkkHIiVKeG5ecnMn1aSYzhoEfAAVIWyYXE0ci3QQSECDQrOYn+/gOIM5U2Lz6WqQfSm8tnBE7rAzxsTFqlvDJ//////AgEAAAAAAAAAGXapFH1x4vjOGUCfuTMSyfAZp97j0U6hiKwMAAAAAAAAABl2qRTlpPBsGiTCqd3/yuRAWq9sIDuj0oisAAAAAA=="
	tx1Source = "AQAAAAH3r1bx2GGVS0wrWlVhnlfTzPJmpx754/9+Wt/v2PhdMAEAAABqRzBEAiAKny96EiIym5fiO8XE3vCCNgBJXir6O8aue+XVu9av8gIgUl9pJF9SVIBmxTGUI1eGiO9GurGLgYVkXb4lh/DJ27NBIQI3r2kZnSB/9c4jCWl2ROrCur6oPq3v2NODhMjK4+h83/////8CAQAAAAAAAAAZdqkUb96slVk3zuwEk+aWBCGEljZyf8CIrA4AAAAAAAAAGXapFJyF4WV/sj1B4UcSSt7233kz6GyyiKwAAAAA"

	// https://whatsonchain.com/tx/88a7c0ed1cb4767cfc8e7434561379eaea21ae78e480cacf4e69284387057c70
	tx2        = "AQAAAAIbSuUDkTFyxeFr2J2rtx01PFucsqHGmXD9TmkOSfl0EAEAAABqRzBEAiAxJ9U+0u2IQ9la0NpJZZ4IbimNyMSr+UZlburh/VyMhgIgXhDzvSw/AcCJA8OWnRONYrHnbJb4VuJ0PPMGnLRpWnVBIQJ5Ili3+6UMih1hVPC0vkpOV7B47+G0eUbAEGl+md3nkf/////IjkyHDWHX4UuJMdlByIj/s2rVjFI2TknC3yD1ZdrezQEAAABrSDBFAiEAxCUx8LUKyrb9H2OzCisQRqwpllvAz0FAmxgNPUuRq+wCIAb7rabUlp3l8WKX+JBaigB2PKGaGiesMbHnfRKHd6FAQSED0h5ymG3g01Sv8d1zegZra3hrwgS+wis5QeEOlXWnqnv/////AhQAAAAAAAAAGXapFOiWQpj8qlBvOebR0fKWV/ecHnLniKwJAAAAAAAAABl2qRTgvT8tXBkZEJgxv61AuOspPAdiG4isAAAAAA=="
	tx2Source1 = "AQAAAAEk7rxBY5UWTwNh9Aqi9VXCbpcV0048BT1OKjIEZa69cQEAAABqRzBEAiAY80ai+e+bl9ELh3G1BiqNxomnz1x9x19AQ7z56chKrQIgZdrUX8Q7Jw6mBQyCod0o4tA8VUTpg6axyBSov3bA15FBIQJkJQ+zNGqqAddYIZ1MVwfK/v4iJPT3juke6lCgVOXXBP////8CAQAAAAAAAAAZdqkU4IQtqp0YqInFfZmqUQ5UkslQv5mIrBAAAAAAAAAAGXapFPhwTZFa19K1WfYbrWwxtg3qxSo3iKwAAAAA"
	tx2Source2 = "AQAAAAHU58f2jMJt3XzGjJEKINLPVzwd2Mr6NDEAq8exla/vIgEAAABrSDBFAiEA3rvUh3L5fGG8nzMdxTW6AoKarzlehm3pHMDDULQ+f0sCIAmo1o/v9WUJD62kTZgsZ3iBYn3AjpkjOG7iWyedxxCxQSEDXI/Xt/qQrisBpMkdoNh/87u8M5DZ3me2n61SqLeP9J3/////AgEAAAAAAAAAGXapFAS8COAvcQwoaykycYzP1nGgyBZEiKwOAAAAAAAAABl2qRRrgpexw82ewTFRyQ0p46lvFHU1poisAAAAAA=="

	// https://whatsonchain.com/tx/85ad54dcbcfa807afee658fe032c64bdef49045fd1be10accb381cf56f32fb61
	txTwoInputsOneSource       = "AQAAAAJQC6Fqb/LLfGobs6yvhVhnIATPVNv/hsXxQAJUnh95LwEAAABqRzBEAiBidLmtZJ7As+5jOMkkHYUS8PfHN38uD45fEE3Ml+UWgwIgAUbOS7ewfSuHja+YanzROzWdASj+qjdYBQhaFURS6+BBIQIjRmq8qJ+aZSw9RUNeG3zmUgojmu/utS6uCuVLCAN7Mf////9QC6Fqb/LLfGobs6yvhVhnIATPVNv/hsXxQAJUnh95LwIAAABqRzBEAiBV8RqlbMPN8pzuSjL+dTiriF1mFT2BRmfoaV9jgIBBNQIgbdVSsTLwNsMZTfZAv2P1JnvKS/8AOS1gpKE1++9HRxFBIQKhVT5dyxO2zlxlv5467kXQmJDYncne0ndKIKV7JY0pS/////8QehJRgwAAAAAZdqkUdLdyQGJltYL3K/qM96Q09BY+PCqIrADKmjsAAAAAGXapFDQD+yX5I8Kq/jEsHl4Rj41KDikIiKwAZc0dAAAAABl2qRRLTGfS7F3X6k5tn5XGgxnbnELPlIisAMLrCwAAAAAZdqkUBk2Kt4LHQ4z2R9zk9BB8b1MixTyIrIDw+gIAAAAAGXapFNkf5iX1ZuuMArDA+Vj68VqqxjOoiKwALTEBAAAAABl2qRSyZBNG4Jz1371skd4PyeyjQl6zTYisAC0xAQAAAAAZdqkUvLivDxZBGS2wd8lxyq6E+Lc638WIrEBLTAAAAAAAGXapFJzi2FpTDUQnATmIe5BmTh0t0kH8iKxAQg8AAAAAABl2qRTz0vxYV0pm8disY/ScP1WNPDWQiYisIKEHAAAAAAAZdqkU/jns4EW0siIyvDCufMRVzLh6ulmIrEANAwAAAAAAGXapFAIc8O+27EndQLPTE+SaX2aibolXiKyghgEAAAAAABl2qRTFHaGdVAgmyYK2M7oxO6Wy54dtx4isUMMAAAAAAAAZdqkUgmvQqSer4kVpiPo2KUkwbCINkkeIrCBOAAAAAAAAGXapFP8UnfvuxI6KooNLnsF0/YPvzyy7iKwAAAAAAAAAAD0AagNkeHMvdXgxM3AzRGZhMTBDSHBoRW9qcGQxWlE9PSxtMTY2LHRiLGExMDAwLGU3MC4xMjMEdGV4dAFCFAYAAAAAAAAZdqkUygwwMJh7r7ptJ85KDq8RZ2153f6IrAAAAAA="
	txTwoInputsOneSourceSource = "AQAAAAHLsc23LvOUkvFjr5pZkZ+iJfWyMan7tKRDdq1lhvUP0AAAAABqRzBEAiB8kOatl0CM4oHG5jKdIQKr3Smy5RYEA9JCGSVuY7MUigIgForNoc0NQi6M/04iPv91n//SYwBwnEp3Z/xLd1B7GilBIQPc7+tXmDGeFkYop4OUydbEuDsmdDlubMsnqEWwWRFlVP////8QPYmoQQAAAAAZdqkUdLdyQGJltYL3K/qM96Q09BY+PCqIrACUNXcAAAAAGXapFAITxKCLLyZm0WyMVRQ1ioKEIDZwiKwAlDV3AAAAABl2qRR0YFvA/NF7Fk3NF7NSjJHq6u5204isAGXNHQAAAAAZdqkUY9pnMfntP7RIhdw7nxUOTtDbk62IrADh9QUAAAAAGXapFOVVjqbmhnEyj12Gcr2+enl098u1iKyA8PoCAAAAABl2qRTWnHHBbtLZcUKHJwZeCDTGuyMaZ4isAC0xAQAAAAAZdqkUOJgc28r9NAikKwGoCQKe767X6LSIrEBLTAAAAAAAGXapFIqwnorF30XabMiyXq3C1AbS+8BTiKxAQg8AAAAAABl2qRQ6sRXhS0zZlS47pd2zKNunaLlkMYisIKEHAAAAAAAZdqkU1n+sMwembspUDwxVteBOcam/FMqIrEANAwAAAAAAGXapFLSWsEFzSPZsofaId13I6CtoSRjjiKyghgEAAAAAABl2qRQCD8wNLFHTTJoDxc1r6uFEchpwUYisIE4AAAAAAAAZdqkUh0BhT5FwJ1vwonWq5NrKCD9tPG2IrBAnAAAAAAAAGXapFBCLoKA2leMnOj9Cmkr27MhAUnsDiKwAAAAAAAAAADoAagNkeHMsdXgxM3AzRGZhMTBDSHBoRW9qcGQxWlE9PSxtMTY2LHRiLGE1MDAsZTczLjkEdGV4dAFCdQoAAAAAAAAZdqkUW+gEsF8m9sRopEv79y9xPdL2AF2IrAAAAAA="
)

func TestConverter(t *testing.T) {
	tests := map[string]struct {
		rawTXBase64   string
		txGetter      *mockTransactionsGetter
		expectedEFHex string
	}{
		"Convert tx with one unsourced input": {
			rawTXBase64: tx1,
			txGetter: newMockTransactionsGetter(t, []string{
				tx1Source,
			}),
			expectedEFHex: "010000000000000000ef01776f9d2c0d80b612ca54ebca1fa3bd38db375756ec7778edddc323569e06dc96010000006b483045022100f83415750880cc9464b752c33215ede7568c45c83cd3ccb841787edd1219d368022065e71b09241c889529e1b979c9cc9f569263386811f0005485b2617134722dd04121020d0ace627fbf80e20ce54d8bcfa5aa41f4a6f2d9c113bac0cf1b1316a96f0c9fffffffff0e000000000000001976a9149c85e1657fb23d41e147124adef6df7933e86cb288ac0201000000000000001976a9147d71e2f8ce19409fb93312c9f019a7dee3d14ea188ac0c000000000000001976a914e5a4f06c1a24c2a9ddffcae4405aaf6c203ba3d288ac00000000",
		},
		"Convert tx with two unsourced inputs": {
			rawTXBase64: tx2,
			txGetter: newMockTransactionsGetter(t, []string{
				tx2Source1,
				tx2Source2,
			}),
			expectedEFHex: "010000000000000000ef021b4ae503913172c5e16bd89dabb71d353c5b9cb2a1c69970fd4e690e49f97410010000006a47304402203127d53ed2ed8843d95ad0da49659e086e298dc8c4abf946656eeae1fd5c8c8602205e10f3bd2c3f01c08903c3969d138d62b1e76c96f856e2743cf3069cb4695a75412102792258b7fba50c8a1d6154f0b4be4a4e57b078efe1b47946c010697e99dde791ffffffff10000000000000001976a914f8704d915ad7d2b559f61bad6c31b60deac52a3788acc88e4c870d61d7e14b8931d941c888ffb36ad58c52364e49c2df20f565dadecd010000006b483045022100c42531f0b50acab6fd1f63b30a2b1046ac29965bc0cf41409b180d3d4b91abec022006fbada6d4969de5f16297f8905a8a00763ca19a1a27ac31b1e77d128777a140412103d21e72986de0d354aff1dd737a066b6b786bc204bec22b3941e10e9575a7aa7bffffffff0e000000000000001976a9146b8297b1c3cd9ec13151c90d29e3a96f147535a688ac0214000000000000001976a914e8964298fcaa506f39e6d1d1f29657f79c1e72e788ac09000000000000001976a914e0bd3f2d5c1919109831bfad40b8eb293c07621b88ac00000000",
		},
		"Convert tx with two inputs from one source": {
			rawTXBase64: txTwoInputsOneSource,
			txGetter: newMockTransactionsGetter(t, []string{
				txTwoInputsOneSourceSource,
			}),
			expectedEFHex: "010000000000000000ef02500ba16a6ff2cb7c6a1bb3acaf8558672004cf54dbff86c5f14002549e1f792f010000006a47304402206274b9ad649ec0b3ee6338c9241d8512f0f7c7377f2e0f8e5f104dcc97e5168302200146ce4bb7b07d2b878daf986a7cd13b359d0128feaa375805085a154452ebe041210223466abca89f9a652c3d45435e1b7ce6520a239aefeeb52eae0ae54b08037b31ffffffff00943577000000001976a9140213c4a08b2f2666d16c8c5514358a828420367088ac500ba16a6ff2cb7c6a1bb3acaf8558672004cf54dbff86c5f14002549e1f792f020000006a473044022055f11aa56cc3cdf29cee4a32fe7538ab885d66153d814667e8695f638080413502206dd552b132f036c3194df640bf63f5267bca4bff00392d60a4a135fbef474711412102a1553e5dcb13b6ce5c65bf9e3aee45d09890d89dc9ded2774a20a57b258d294bffffffff00943577000000001976a91474605bc0fcd17b164dcd17b3528c91eaeaee76d388ac107a125183000000001976a91474b772406265b582f72bfa8cf7a434f4163e3c2a88ac00ca9a3b000000001976a9143403fb25f923c2aafe312c1e5e118f8d4a0e290888ac0065cd1d000000001976a9144b4c67d2ec5dd7ea4e6d9f95c68319db9c42cf9488ac00c2eb0b000000001976a914064d8ab782c7438cf647dce4f4107c6f5322c53c88ac80f0fa02000000001976a914d91fe625f566eb8c02b0c0f958faf15aaac633a888ac002d3101000000001976a914b2641346e09cf5dfbd6c91de0fc9eca3425eb34d88ac002d3101000000001976a914bcb8af0f1641192db077c971caae84f8b73adfc588ac404b4c00000000001976a9149ce2d85a530d44270139887b90664e1d2dd241fc88ac40420f00000000001976a914f3d2fc58574a66f1d8ac63f49c3f558d3c35908988ac20a10700000000001976a914fe39ece045b4b22232bc30ae7cc455ccb87aba5988ac400d0300000000001976a914021cf0efb6ec49dd40b3d313e49a5f66a26e895788aca0860100000000001976a914c51da19d540826c982b633ba313ba5b2e7876dc788ac50c30000000000001976a914826bd0a927abe2456988fa362949306c220d924788ac204e0000000000001976a914ff149dfbeec48e8aa2834b9ec174fd83efcf2cbb88ac00000000000000003d006a036478732f757831337033446661313043487068456f6a7064315a513d3d2c6d3136362c74622c61313030302c6537302e3132330474657874014214060000000000001976a914ca0c3030987bafba6d27ce4a0eaf11676d79ddfe88ac00000000",
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			tx := txFromBase64(t, test.rawTXBase64)

			converter := ef.NewConverter(test.txGetter)
			efHex, err := converter.Convert(context.Background(), tx)
			require.NoError(t, err)
			require.Equal(t, test.expectedEFHex, efHex)

			// additionally check if converter can convert already-in-ef-hex tx
			// without the need to fetch source transactions
			tx, err = sdk.NewTransactionFromHex(efHex)
			require.NoError(t, err)
			converter = ef.NewConverter(newMockTransactionsGetter(t, []string{}))
			efHexRegenerated, err := converter.Convert(context.Background(), tx)
			require.NoError(t, err)
			require.Equal(t, efHex, efHexRegenerated)
		})
	}
}

func TestConverterErrorCases(t *testing.T) {
	tests := map[string]struct {
		rawTXBase64 string
		txGetter    *mockTransactionsGetter
		expectErr   error
	}{
		"No source tx provided by TransactionGetter": {
			rawTXBase64: tx1,
			txGetter:    newMockTransactionsGetter(t, []string{}).WithOnMissingBehavior(onMissingTxSkip),
			expectErr:   ef.ErrGetTransactions,
		},
		"Not every source tx provided by TransactionGetter": {
			rawTXBase64: tx2,
			txGetter: newMockTransactionsGetter(t, []string{
				tx2Source1,
			}).WithOnMissingBehavior(onMissingTxSkip),
			expectErr: ef.ErrGetTransactions,
		},
		"TransactionGetter error on missing transaction": {
			rawTXBase64: tx1,
			txGetter:    newMockTransactionsGetter(t, []string{}).WithOnMissingBehavior(onMissingTxReturnError),
			expectErr:   ef.ErrGetTransactions,
		},
		"Nil transaction returned by TransactionGetter": {
			rawTXBase64: tx1,
			txGetter:    newMockTransactionsGetter(t, []string{}).WithOnMissingBehavior(onMissingTxAddNil),
			expectErr:   ef.ErrGetTransactions,
		},
		"TransactionGetter returned more transactions than requested": {
			rawTXBase64: tx1,
			txGetter:    newMockTransactionsGetter(t, []string{tx1Source, tx2Source1}).WithReturnAll(true),
			expectErr:   ef.ErrGetTransactions,
		},
		"TransactionGetter not requested transactions but with correct length": {
			rawTXBase64: tx1,
			txGetter:    newMockTransactionsGetter(t, []string{tx2Source1}).WithReturnAll(true),
			expectErr:   ef.ErrGetTransactions,
		},
		"TransactionGetter duplicated transaction": {
			rawTXBase64: tx2,
			txGetter:    newMockTransactionsGetter(t, []string{tx2Source1, tx2Source1}).WithReturnAll(true),
			expectErr:   ef.ErrGetTransactions,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			tx := txFromBase64(t, test.rawTXBase64)

			converter := ef.NewConverter(test.txGetter)
			efHex, err := converter.Convert(context.Background(), tx)
			require.ErrorIs(t, err, test.expectErr)
			require.Empty(t, efHex)
		})
	}
}

func txFromBase64(t *testing.T, rawTXBase64 string) *sdk.Transaction {
	d, err := base64.StdEncoding.DecodeString(rawTXBase64)
	require.NoError(t, err)
	tx, err := sdk.NewTransactionFromBytes(d)
	require.NoError(t, err)
	return tx
}
