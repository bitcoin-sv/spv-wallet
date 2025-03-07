package contacts_test

import (
	"fmt"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	"github.com/bitcoin-sv/spv-wallet/actions/testabilities/apierror"
	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
)

func TestRejectContact(t *testing.T) {
	// given:
	givenForAllTests := testabilities.Given(t)
	cleanup := givenForAllTests.StartedSPVWalletWithConfiguration(
		testengine.WithV2(),
	)
	defer cleanup()

	t.Run("Reject contact", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		given.User(fixtures.Sender).HasAwaitingContactTo(fixtures.RecipientInternal)
		client := given.HttpClient().ForGivenUser(fixtures.Sender)

		// when:
		res, _ := client.R().
			Delete(fmt.Sprintf("/api/v2/invitations/%s", fixtures.RecipientInternal.DefaultPaymail().String()))

		// then:
		then.Response(res).IsOK()
	})

	t.Run("Reject already rejected contact", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		given.User(fixtures.Sender).HasAwaitingContactTo(fixtures.RecipientInternal)
		client := given.HttpClient().ForGivenUser(fixtures.Sender)

		// when:
		res, _ := client.R().
			Delete(fmt.Sprintf("/api/v2/invitations/%s", fixtures.RecipientInternal.DefaultPaymail().String()))

		// then:
		then.Response(res).IsOK()

		// and:
		// when:
		res, _ = client.R().
			Delete(fmt.Sprintf("/api/v2/invitations/%s", fixtures.RecipientInternal.DefaultPaymail().String()))

		// then:
		then.Response(res).
			HasStatus(404).
			WithJSONf(apierror.ExpectedJSON("error-contact-not-found", "contact not found"))
	})

	t.Run("Contact in wrong status", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		given.User(fixtures.Sender).HasContactTo(fixtures.RecipientInternal)
		client := given.HttpClient().ForGivenUser(fixtures.Sender)

		// when:
		res, _ := client.R().
			Delete(fmt.Sprintf("/api/v2/invitations/%s", fixtures.RecipientInternal.DefaultPaymail().String()))

		// then:
		then.Response(res).
			HasStatus(400).
			WithJSONf(apierror.ExpectedJSON("error-contact-wrong-status", "contact is in wrong status"))
	})

	t.Run("Reject contact with admin xpub", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		given.User(fixtures.Sender).HasContactTo(fixtures.RecipientInternal)

		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.R().
			Delete(fmt.Sprintf("/api/v2/invitations/%s", fixtures.RecipientInternal.DefaultPaymail().String()))

		// then:
		then.Response(res).
			HasStatus(401).
			WithJSONf(apierror.ExpectedJSON("error-admin-auth-on-user-endpoint", "cannot call user's endpoints with admin authorization"))
	})

	t.Run("No contact to reject", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForGivenUser(fixtures.Sender)

		// when:
		res, _ := client.R().
			Delete(fmt.Sprintf("/api/v2/invitations/%s", fixtures.RecipientExternal.DefaultPaymail().String()))

		// then:
		then.Response(res).
			HasStatus(404).
			WithJSONf(apierror.ExpectedJSON("error-contact-not-found", "contact not found"))
	})
}
