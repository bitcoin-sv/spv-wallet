package contacts_test

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
)

func TestGetContact(t *testing.T) {
	// given:
	givenForAllTests := testabilities.Given(t)
	cleanup := givenForAllTests.StartedSPVWalletWithConfiguration(
		testengine.WithV2(),
	)
	defer cleanup()

	t.Run("No contact to get", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForGivenUser(fixtures.Sender)

		// when:
		res, _ := client.R().Get("/api/v2/contacts/" + fixtures.UserWithMorePaymails.DefaultPaymail().String())

		// then:
		then.Response(res).
			HasStatus(404).
			WithJSONMatching(`{
				"code": "{{ .code }}",
				"message": "{{ .message }}"
			}`, map[string]any{
				"code":    spverrors.ErrContactNotFound.Code,
				"message": spverrors.ErrContactNotFound.Message,
			})
	})

	t.Run("Get contact", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		given.User(fixtures.Sender).HasContactTo(fixtures.RecipientInternal)
		given.User(fixtures.RecipientInternal).HasContactTo(fixtures.Sender)
		c3 := given.User(fixtures.Sender).HasContactTo(fixtures.UserWithMorePaymails)
		client := given.HttpClient().ForGivenUser(fixtures.Sender)

		// when:
		res, _ := client.R().Get("/api/v2/contacts/" + fixtures.UserWithMorePaymails.DefaultPaymail().String())

		// then:
		then.Response(res).
			HasStatus(200).
			WithJSONMatching(`
			{
				"id": "{{ matchNumber }}",
				"createdAt": "{{ matchTimestamp }}",
				"updatedAt": "{{ matchTimestamp }}",
				"fullName": "{{ .fullName }}",
				"paymail": "{{ .paymail }}",
				"pubKey": "{{ matchHexWithLength 66 }}",
				"status": "{{ .status }}"
			}`, map[string]any{
				"fullName": c3.FullName,
				"paymail":  c3.Paymail,
				"status":   c3.Status,
			})
	})

	t.Run("Get contact with admin xpub", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		given.User(fixtures.Sender).HasContactTo(fixtures.RecipientInternal)
		given.User(fixtures.RecipientInternal).HasContactTo(fixtures.Sender)
		given.User(fixtures.Sender).HasContactTo(fixtures.UserWithMorePaymails)
		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.R().Get("/api/v2/contacts/" + fixtures.UserWithMorePaymails.DefaultPaymail().String())

		// then:
		then.Response(res).
			HasStatus(401).
			WithJSONMatching(`{
				"code": "{{ .code }}",
				"message": "{{ .message }}"
			}`, map[string]any{
				"code":    spverrors.ErrAdminAuthOnUserEndpoint.Code,
				"message": spverrors.ErrAdminAuthOnUserEndpoint.Message,
			})
	})
}
