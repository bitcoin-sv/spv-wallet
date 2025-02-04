package testabilities

import (
	"maps"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
)

// ExpectedCustomInstructions represents json array with custom instructions of UTXO created with Faucet().TopUp() method
const ExpectedCustomInstructions = `[
		{
			"type": "sign",
			"instruction": "P2PKH"
		}
	]`

type TransactionsEndpointFixture interface {
	testabilities.SPVWalletApplicationFixture
	OutlineResponseContext(additionalParams map[string]any) map[string]any
}

func Given(t testing.TB) TransactionsEndpointFixture {
	return &fixture{
		SPVWalletApplicationFixture: testabilities.Given(t),
	}
}

type fixture struct {
	testabilities.SPVWalletApplicationFixture
}

func (f *fixture) OutlineResponseContext(paramsOverride map[string]any) map[string]any {
	params := make(map[string]any)
	params["ReceiverPaymail"] = fixtures.RecipientExternal.DefaultPaymail()
	params["SenderPaymail"] = fixtures.Sender.DefaultPaymail()
	params["CustomInstructions"] = ExpectedCustomInstructions

	maps.Copy(params, paramsOverride)

	return params
}
