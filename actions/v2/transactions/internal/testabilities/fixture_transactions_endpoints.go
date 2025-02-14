package testabilities

import (
	"maps"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/paymailmock"
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
	OutlineResponseContext(format string, additionalParams map[string]any) map[string]any
}

func Given(t testing.TB) TransactionsEndpointFixture {
	return &fixture{
		SPVWalletApplicationFixture: testabilities.Given(t),
	}
}

type fixture struct {
	testabilities.SPVWalletApplicationFixture
}

func (f *fixture) OutlineResponseContext(format string, paramsOverride map[string]any) map[string]any {
	params := map[string]any{
		"Format":             format,
		"ReceiverPaymail":    fixtures.RecipientExternal.DefaultPaymail(),
		"SenderPaymail":      fixtures.Sender.DefaultPaymail(),
		"CustomInstructions": ExpectedCustomInstructions,
		"Reference":          paymailmock.MockedReferenceID,
	}

	maps.Copy(params, paramsOverride)

	return params
}
