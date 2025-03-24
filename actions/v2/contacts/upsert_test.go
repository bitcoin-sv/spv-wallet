package contacts_test

import (
	"fmt"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	"github.com/bitcoin-sv/spv-wallet/actions/testabilities/apierror"
	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"github.com/bitcoin-sv/spv-wallet/engine/v2/contacts/contactsmodels"
)

func TestUpsertContact(t *testing.T) {
	// given:
	givenForAllTests := testabilities.Given(t)
	cleanup := givenForAllTests.StartedSPVWalletWithConfiguration(
		testengine.WithV2(),
	)
	defer cleanup()

	t.Run("Create contact", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)

		client := given.HttpClient().ForGivenUser(fixtures.Sender)

		// when:
		res, _ := client.R().
			SetBody(map[string]any{
				"requesterPaymail": fixtures.Sender.DefaultPaymail().String(),
				"fullName":         fixtures.RecipientInternal.DefaultPaymail().PublicName(),
			}).
			Put(fmt.Sprintf("/api/v2/contacts/%s", fixtures.RecipientInternal.DefaultPaymail().String()))

		// then:
		then.Response(res).
			HasStatus(200).
			WithJSONMatching(`{
				"id": "{{ matchNumber }}",
				"createdAt": "{{ matchTimestamp }}",
				"updatedAt": "{{ matchTimestamp }}",
				"fullName": "{{ .fullName }}",
				"paymail": "{{ .paymail }}",
				"pubKey": "{{ matchHexWithLength 66 }}",
				"status": "{{ .status }}"
			}`, map[string]any{
				"fullName": fixtures.RecipientInternal.DefaultPaymail().PublicName(),
				"paymail":  fixtures.RecipientInternal.DefaultPaymail().String(),
				"status":   contactsmodels.ContactNotConfirmed,
			})
	})

	t.Run("Upsert contact with admin xpub", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)

		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.R().
			SetBody(map[string]any{
				"requesterPaymail": fixtures.Sender.DefaultPaymail().String(),
				"fullName":         fixtures.RecipientInternal.DefaultPaymail().PublicName(),
			}).
			Put(fmt.Sprintf("/api/v2/contacts/%s", fixtures.RecipientInternal.DefaultPaymail().String()))

		// then:
		then.Response(res).
			HasStatus(401).
			WithJSONf(apierror.ExpectedJSON("error-admin-auth-on-user-endpoint", "cannot call user's endpoints with admin authorization"))
	})

	t.Run("Create contact with not found requester paymail", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)

		client := given.HttpClient().ForGivenUser(fixtures.Sender)

		// when:
		res, _ := client.R().
			SetBody(map[string]any{
				"requesterPaymail": fixtures.RecipientExternal.DefaultPaymail().String(),
				"fullName":         fixtures.RecipientInternal.DefaultPaymail().PublicName(),
			}).
			Put(fmt.Sprintf("/api/v2/contacts/%s", fixtures.RecipientInternal.DefaultPaymail().String()))

		// then:
		then.Response(res).
			HasStatus(404).
			WithJSONf(apierror.ExpectedJSON("error-paymail-not-found", "paymail not found"))
	})

	t.Run("Create contact with mismatching requester paymail", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)

		client := given.HttpClient().ForGivenUser(fixtures.Sender)

		// when:
		res, _ := client.R().
			SetBody(map[string]any{
				"requesterPaymail": fixtures.RecipientInternal.DefaultPaymail().String(),
				"fullName":         fixtures.RecipientInternal.DefaultPaymail().PublicName(),
			}).
			Put(fmt.Sprintf("/api/v2/contacts/%s", fixtures.RecipientInternal.DefaultPaymail().String()))

		// then:
		then.Response(res).
			HasStatus(400).
			WithJSONf(apierror.ExpectedJSON("error-paymail-user-do-not-own", "user do not own paymail"))
	})

	t.Run("Create contact without creator paymail", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		given.Paymail().ExternalPaymailHost().WillRespondWithBasicCapabilities()

		client := given.HttpClient().ForGivenUser(fixtures.Sender)

		// when:
		res, _ := client.R().
			SetBody(map[string]any{
				"fullName": fixtures.RecipientInternal.DefaultPaymail().PublicName(),
			}).
			Put(fmt.Sprintf("/api/v2/contacts/%s", fixtures.RecipientInternal.DefaultPaymail().String()))

		// then:
		then.Response(res).
			HasStatus(404).
			WithJSONf(apierror.ExpectedJSON("error-paymail-not-found", "paymail not found"))
	})

	t.Run("Create contact without full name", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		given.Paymail().ExternalPaymailHost().WillRespondWithBasicCapabilities()

		client := given.HttpClient().ForGivenUser(fixtures.Sender)

		// when:
		res, _ := client.R().
			SetBody(map[string]any{
				"requesterPaymail": fixtures.Sender.DefaultPaymail().String(),
			}).
			Put(fmt.Sprintf("/api/v2/contacts/%s", fixtures.RecipientInternal.DefaultPaymail().String()))

		// then:
		then.Response(res).
			HasStatus(400).
			WithJSONf(apierror.ExpectedJSON("error-contact-full-name-required", "full name is required"))
	})
}
