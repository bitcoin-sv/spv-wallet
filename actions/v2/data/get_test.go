package data_test

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	"github.com/bitcoin-sv/spv-wallet/actions/testabilities/apierror"
	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
)

// NOTE: More complex and real-world test case can be found in outlines_record_test.go
func TestGetData(t *testing.T) {
	// given:
	given, then := testabilities.New(t)
	cleanup := given.StartedSPVWalletWithConfiguration(
		testengine.WithV2(),
	)
	defer cleanup()

	// and:
	dataToStore := "hello world"

	// and:
	_, dataID := given.Faucet(fixtures.Sender).StoreData(dataToStore)

	// and:
	client := given.HttpClient().ForGivenUser(fixtures.Sender)

	// when:
	res, _ := client.R().
		Get("/api/v2/data/" + dataID)

	// then:
	then.Response(res).
		IsOK().WithJSONMatching(`{
				"id": "{{ .outpoint }}",
				"blob": "{{ .blob }}"
			}`, map[string]any{
		"outpoint": dataID,
		"blob":     dataToStore,
	})
}

func TestErrorCases(t *testing.T) {
	// given:
	givenForAllTests := testabilities.Given(t)
	cleanup := givenForAllTests.StartedSPVWalletWithConfiguration(
		testengine.WithV2(),
	)
	defer cleanup()

	// and
	mockTx := givenForAllTests.Tx().WithInput(1).WithP2PKHOutput(1)
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

	t.Run("try to get data with wrong id", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)

		// and:
		client := given.HttpClient().ForUser()

		// and:
		wrongID := "wrong_id" // doesn't match the outpoint format "<txID>-<vout>"

		// when:
		res, _ := client.R().Get("/api/v2/data/" + wrongID)

		// then:
		then.Response(res).HasStatus(400).WithJSONf(
			apierror.ExpectedJSON("error-invalid-data-id", "invalid data id"),
		)
	})
}
