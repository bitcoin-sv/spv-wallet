package contacts_test

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	"github.com/bitcoin-sv/spv-wallet/actions/testabilities/apierror"
	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
)

func TestConfirmContact(t *testing.T) {
	// given:
	givenForAllTests := testabilities.Given(t)
	cleanup := givenForAllTests.StartedSPVWalletWithConfiguration(
		testengine.WithV2(),
	)
	defer cleanup()

	t.Run("No side has contact", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)

		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.R().
			SetBody(map[string]any{
				"paymailA": fixtures.Sender.DefaultPaymail().String(),
				"paymailB": fixtures.RecipientInternal.DefaultPaymail().String(),
			}).
			Post("/api/v2/admin/contacts/confirmations")

		// then:
		then.Response(res).HasStatus(500).
			WithJSONf(apierror.ExpectedJSON("error-contact-getting-contact-failed", "getting contact failed"))
	})

	t.Run("Only one side has contact", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		given.User(fixtures.Sender).HasContactTo(fixtures.RecipientInternal)

		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.R().
			SetBody(map[string]any{
				"paymailA": fixtures.Sender.DefaultPaymail().String(),
				"paymailB": fixtures.RecipientInternal.DefaultPaymail().String(),
			}).
			Post("/api/v2/admin/contacts/confirmations")

		// then:
		then.Response(res).HasStatus(500).
			WithJSONf(apierror.ExpectedJSON("error-contact-getting-contact-failed", "getting contact failed"))
	})

	t.Run("Confirm contact", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		given.User(fixtures.Sender).HasContactTo(fixtures.RecipientInternal)
		given.User(fixtures.RecipientInternal).HasContactTo(fixtures.Sender)

		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.R().
			SetBody(map[string]any{
				"paymailA": fixtures.Sender.DefaultPaymail().String(),
				"paymailB": fixtures.RecipientInternal.DefaultPaymail().String(),
			}).
			Post("/api/v2/admin/contacts/confirmations")

		// then:
		then.Response(res).IsOK()
	})

	t.Run("Confirm already confirmed contact", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		given.User(fixtures.Sender).HasContactTo(fixtures.RecipientInternal)
		given.User(fixtures.RecipientInternal).HasContactTo(fixtures.Sender)

		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.R().
			SetBody(map[string]any{
				"paymailA": fixtures.Sender.DefaultPaymail().String(),
				"paymailB": fixtures.RecipientInternal.DefaultPaymail().String(),
			}).
			Post("/api/v2/admin/contacts/confirmations")

		// then:
		then.Response(res).IsOK()

		// when:
		res, _ = client.R().
			SetBody(map[string]any{
				"paymailA": fixtures.Sender.DefaultPaymail().String(),
				"paymailB": fixtures.RecipientInternal.DefaultPaymail().String(),
			}).
			Post("/api/v2/admin/contacts/confirmations")

		// then:
		then.Response(res).IsOK()
	})

	t.Run("Confirm contact with user xpub", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)

		client := given.HttpClient().ForUser()

		// when:
		res, _ := client.R().
			SetBody(map[string]any{
				"paymailA": fixtures.RecipientExternal.DefaultPaymail().String(),
				"paymailB": fixtures.Sender.DefaultPaymail().String(),
			}).
			Post("/api/v2/admin/contacts/confirmations")

		// then:
		then.Response(res).
			HasStatus(401).
			WithJSONf(apierror.ExpectedJSON("error-unauthorized-xpub-not-an-admin-key", "xpub provided is not an admin key"))
	})

}
