package transactions_test

import (
	"fmt"
	"github.com/bitcoin-sv/spv-wallet/actions/testabilities/apierror"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
)

//NOTE: Standard Flow Case is tested in outlines_record_test.go

func TestGetTransactionDataErrorCases(t *testing.T) {
	t.Run("no data under outpoint", func(t *testing.T) {
		// given:
		given, then := testabilities.New(t)
		cleanup := given.StartedSPVWallet()
		defer cleanup()
		client := given.HttpClient().ForGivenUser(fixtures.Sender)

		// when:
		res, _ := client.R().Get(fmt.Sprintf("%s/%s/0", getTransactionDataURL, fixtures.GivenTX(t).ID()))

		// then:
		then.Response(res).
			HasStatus(404).
			WithJSONf(apierror.ExpectedJSON("error-transaction-data-outpoint-not-found", "data outpoint not found"))
	})

	t.Run("try to return paymails info for admin", func(t *testing.T) {
		// given:
		given, then := testabilities.New(t)
		cleanup := given.StartedSPVWallet()
		defer cleanup()
		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.R().Get(fmt.Sprintf("%s/%s/0", getTransactionDataURL, fixtures.GivenTX(t).ID()))

		// then:
		then.Response(res).IsUnauthorizedForAdmin()
	})

	t.Run("return xpub info for anonymous", func(t *testing.T) {
		// given:
		given, then := testabilities.New(t)
		cleanup := given.StartedSPVWallet()
		defer cleanup()
		client := given.HttpClient().ForAnonymous()

		// when:
		res, _ := client.R().Get(fmt.Sprintf("%s/%s/0", getTransactionDataURL, fixtures.GivenTX(t).ID()))

		// then:
		then.Response(res).IsUnauthorized()
	})
}
