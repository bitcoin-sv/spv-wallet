package contacts_test

import (
	"fmt"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
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

	t.Run("Confirm contact", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		given.User(fixtures.Sender).HasContactTo(fixtures.RecipientInternal)
		client := given.HttpClient().ForGivenUser(fixtures.Sender)

		// when:
		res, _ := client.R().
			Post(fmt.Sprintf("/api/v2/contacts/%s/confirmation", fixtures.RecipientInternal.DefaultPaymail().String()))

		// then:
		then.Response(res).IsOK()
	})

	t.Run("Confirm already confirmed contact", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		given.User(fixtures.Sender).HasContactTo(fixtures.RecipientInternal)
		client := given.HttpClient().ForGivenUser(fixtures.Sender)

		// when:
		res, _ := client.R().
			Post(fmt.Sprintf("/api/v2/contacts/%s/confirmation", fixtures.RecipientInternal.DefaultPaymail().String()))

		// then:
		then.Response(res).IsOK()

		// and:
		// when:
		res, _ = client.R().
			Post(fmt.Sprintf("/api/v2/contacts/%s/confirmation", fixtures.RecipientInternal.DefaultPaymail().String()))

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

	t.Run("Contact in wrong status", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		given.User(fixtures.Sender).HasAwaitingContactTo(fixtures.RecipientInternal)
		client := given.HttpClient().ForGivenUser(fixtures.Sender)

		// when:
		res, _ := client.R().
			Post(fmt.Sprintf("/api/v2/contacts/%s/confirmation", fixtures.RecipientInternal.DefaultPaymail().String()))

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

	t.Run("Confirm contact with admin xpub", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		given.User(fixtures.Sender).HasContactTo(fixtures.RecipientInternal)

		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.R().
			Post(fmt.Sprintf("/api/v2/contacts/%s/confirmation", fixtures.RecipientInternal.DefaultPaymail().String()))

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

	t.Run("No contact to confirm", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForGivenUser(fixtures.Sender)

		// when:
		res, _ := client.R().
			Post(fmt.Sprintf("/api/v2/contacts/%s/confirmation", fixtures.RecipientExternal.DefaultPaymail().String()))

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
