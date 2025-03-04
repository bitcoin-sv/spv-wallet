package contacts_test

import (
	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"testing"
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
			Delete("/api/v2/invitations/" + fixtures.RecipientInternal.DefaultPaymail().String())

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
			Delete("/api/v2/invitations/" + fixtures.RecipientInternal.DefaultPaymail().String())

		// then:
		then.Response(res).IsOK()

		// and:
		// when:
		res, _ = client.R().
			Delete("/api/v2/invitations/" + fixtures.RecipientInternal.DefaultPaymail().String())

		// then:
		then.Response(res).
			HasStatus(spverrors.ErrContactNotFound.StatusCode).
			WithJSONMatching(`{
				"code": "{{ .code }}",
				"message": "{{ .message }}"
			}`, map[string]any{
				"code":    spverrors.ErrContactNotFound.Code,
				"message": spverrors.ErrContactNotFound.Message,
			})
	})

	t.Run("Contact in wrong status", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		given.User(fixtures.Sender).HasContactTo(fixtures.RecipientInternal)
		client := given.HttpClient().ForGivenUser(fixtures.Sender)

		// when:
		res, _ := client.R().
			Delete("/api/v2/invitations/" + fixtures.RecipientInternal.DefaultPaymail().String())

		// then:
		then.Response(res).
			HasStatus(spverrors.ErrContactInWrongStatus.StatusCode).
			WithJSONMatching(`{
				"code": "{{ .code }}",
				"message": "{{ .message }}"
			}`, map[string]any{
				"code":    spverrors.ErrContactInWrongStatus.Code,
				"message": spverrors.ErrContactInWrongStatus.Message,
			})
	})

	t.Run("Reject contact with admin xpub", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		given.User(fixtures.Sender).HasContactTo(fixtures.RecipientInternal)

		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.R().
			Delete("/api/v2/invitations/" + fixtures.RecipientInternal.DefaultPaymail().String())

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

	t.Run("No contact to reject", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForGivenUser(fixtures.Sender)

		// when:
		res, _ := client.R().
			Post("/api/v2/contacts/" + fixtures.RecipientExternal.DefaultPaymail().String() + "/confirmation")

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
}
