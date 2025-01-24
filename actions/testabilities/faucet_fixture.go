package testabilities

import (
	"context"
	"fmt"
	"github.com/bitcoin-sv/go-sdk/script"
	chainmodels "github.com/bitcoin-sv/spv-wallet/engine/chain/models"
	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/stretchr/testify/assert"
	"testing"
)

type faucetFixture struct {
	engineWithConfig testengine.EngineWithConfig
	httpClient       SPVWalletHttpClientFixture
	user             fixtures.User
	t                testing.TB
	assert           *assert.Assertions
	arc              ARCFixture
	bhs              BlockHeadersServiceFixture
}

func (f *faucetFixture) TopUp(satoshis bsv.Satoshis) (fixtures.GivenTXSpec, bsv.CustomInstructions) {
	f.t.Helper()

	anonymousClient := f.httpClient.ForAnonymous()
	recipientClient := f.httpClient.ForGivenUser(f.user)

	var destination struct {
		Outputs   []map[string]any
		Reference string
	}

	res, err := anonymousClient.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]any{
			"satoshis": satoshis,
		}).
		SetResult(&destination).
		Post(
			fmt.Sprintf(
				"https://example.com/v1/bsvalias/p2p-payment-destination/%s",
				f.user.DefaultPaymail(),
			),
		)
	f.assert.NoError(err)
	f.assert.Equal(200, res.StatusCode())

	lockingScript, err := script.NewFromHex(destination.Outputs[0]["script"].(string))
	f.assert.NoError(err)
	f.assert.True(lockingScript.IsP2PKH())

	txSpec := fixtures.GivenTX(f.t).
		WithSender(fixtures.ExternalFaucet).
		WithInput(uint64(satoshis+1)).
		WithOutputScript(uint64(satoshis), lockingScript)

	f.arc.WillRespondForBroadcast(200, &chainmodels.TXInfo{
		TxID:     txSpec.ID(),
		TXStatus: chainmodels.SeenOnNetwork,
	})

	f.bhs.WillRespondForMerkleRootsVerify(200, &chainmodels.MerkleRootsConfirmations{
		ConfirmationState: chainmodels.MRConfirmed,
	})

	res, err = anonymousClient.R().
		SetHeader("Content-Type", "application/json").
		SetBody(map[string]any{
			"beef":      txSpec.BEEF(),
			"reference": destination.Reference,
			"metadata": map[string]any{
				"note":   "top-up",
				"sender": fixtures.ExternalFaucet.DefaultPaymail(),
			},
		}).
		Post(
			fmt.Sprintf(
				"https://example.com/v1/bsvalias/beef/%s",
				f.user.DefaultPaymail(),
			),
		)

	f.assert.NoError(err)
	f.assert.Equal(200, res.StatusCode())

	var userInfo struct {
		CurrentBalance bsv.Satoshis `json:"currentBalance"`
	}
	res, err = recipientClient.R().SetResult(&userInfo).Get("/api/v2/users/current")
	f.assert.NoError(err)
	f.assert.Equal(200, res.StatusCode())

	f.assert.GreaterOrEqual(userInfo.CurrentBalance, satoshis)

	address, err := lockingScript.Address()
	f.assert.NoError(err)

	addresses, err := f.engineWithConfig.Engine.AddressesService().FindByStringAddresses(context.Background(), func(yield func(string) bool) {
		yield(address.AddressString)
	})
	f.assert.NoError(err)
	f.assert.Len(addresses, 1)

	return txSpec, addresses[0].CustomInstructions
}
