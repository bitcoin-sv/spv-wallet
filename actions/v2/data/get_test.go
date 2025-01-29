package data_test

import (
	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"testing"
)

func TestErrorCases(t *testing.T) {
	// given:
	givenForAllTests := testabilities.Given(t)
	cleanup := givenForAllTests.StartedSPVWalletWithConfiguration(
		testengine.WithV2(),
	)
	defer cleanup()

	// and
	mockTx := fixtures.GivenTX(t).WithInput(1).WithP2PKHOutput(1)
	mockOutpoint := bsv.Outpoint{
		TxID: mockTx.ID(),
		Vout: 0,
	}

	t.Run("try to get data as admin", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)

		// and:
		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.R().Get("/api/v2/data/" + mockOutpoint.String())

		// then:
		then.Response(res).IsUnauthorizedForAdmin()
	})

	t.Run("try to get data as anonymous", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)

		// and:
		client := given.HttpClient().ForAnonymous()

		// when:
		res, _ := client.R().Get("/api/v2/data/" + mockOutpoint.String())

		// then:
		then.Response(res).IsUnauthorized()
	})
}
