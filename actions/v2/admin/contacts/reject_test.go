package contacts_test

import (
	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	"strconv"
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
		contact := given.User(fixtures.Sender).HasContactTo(fixtures.RecipientInternal)
		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.R().
			Delete("/api/v2/admin/invitations/" + strconv.Itoa(int(contact.ID)))

		// then:
		then.Response(res).IsOK()
	})

	t.Run("Reject already rejected contact", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		contact := given.User(fixtures.Sender).HasContactTo(fixtures.RecipientInternal)
		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.R().
			Delete("/api/v2/admin/invitations/" + strconv.Itoa(int(contact.ID)))

		// then:
		then.Response(res).IsOK()

		// and:
		// when:
		res, _ = client.R().
			Delete("/api/v2/admin/invitations/" + strconv.Itoa(int(contact.ID)))

		// then:
		then.Response(res).IsOK()
	})

	t.Run("Reject contact with user xpub", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		contact := given.User(fixtures.Sender).HasContactTo(fixtures.RecipientInternal)

		client := given.HttpClient().ForUser()

		// when:
		res, _ := client.R().
			Delete("/api/v2/admin/invitations/" + strconv.Itoa(int(contact.ID)))

		// then:
		then.Response(res).
			HasStatus(401).
			WithJSONMatching(`{
				"code": "{{ .code }}",
				"message": "{{ .message }}"
			}`, map[string]any{
				"code":    spverrors.ErrNotAnAdminKey.Code,
				"message": spverrors.ErrNotAnAdminKey.Message,
			})
	})

	t.Run("No contact to reject", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.R().
			Delete("/api/v2/admin/invitations/99999")

		// then:
		then.Response(res).
			HasStatus(500).
			WithJSONMatching(`{
				"code": "{{ .code }}",
				"message": "{{ .message }}"
			}`, map[string]any{
				"code":    spverrors.ErrUpdateContactStatus.Code,
				"message": spverrors.ErrUpdateContactStatus.Message,
			})
	})
}
