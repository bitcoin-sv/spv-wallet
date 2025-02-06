package admin_test

import (
	"net/http"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
	assert "github.com/stretchr/testify/require"
)

func TestGetPaymails(t *testing.T) {
	// given:
	givenForAllTests := testabilities.Given(t)
	cleanup := givenForAllTests.StartedSPVWallet()
	defer cleanup()

	// and
	userPaymail := fixtures.Sender.DefaultPaymail()

	var testState struct {
		defaultPaymailID string
	}

	t.Run("get paymails for selected user as admin", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.R().Get("/api/v1/admin/paymails?alias=" + userPaymail.Alias())

		// then:
		then.Response(res).
			IsOK().
			WithJSONMatching(`{
			 "content": [
				{
				  "address": "{{.Address}}",
				  "alias": "{{.Alias}}",
				  "avatar": "{{ matchURL | orEmpty }}",
				  "createdAt": "{{ matchTimestamp }}",
				  "deletedAt": null,
				  "domain": "{{.Domain}}",
				  "id": "{{ matchID64 }}",
				  "metadata": "*",
				  "publicName": "{{.PublicName}}",
				  "updatedAt": "{{ matchTimestamp }}",
				  "xpubId": "{{.XPubID}}"
				}
			 ],
			 "page": {
				"number": 1,
				"size": 50,
				"totalElements": 1,
				"totalPages": 1
			 }
			}`, map[string]any{
				"Address":    userPaymail,
				"PublicName": userPaymail.PublicName(),
				"Alias":      userPaymail.Alias(),
				"XPubID":     fixtures.Sender.XPubID(),
				"Domain":     fixtures.PaymailDomain,
			})

		// update:
		getter := then.Response(res).JSONValue()
		testState.defaultPaymailID = getter.GetString("content[0]/id")
	})

	t.Run("get single paymail as admin", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.R().Get("/api/v1/admin/paymails/" + testState.defaultPaymailID)

		// then:
		then.Response(res).
			IsOK().
			WithJSONMatching(`{
				  "address": "{{.Address}}",
				  "alias": "{{.Alias}}",
				  "avatar": "{{ matchURL | orEmpty }}",
				  "createdAt": "{{ matchTimestamp }}",
				  "deletedAt": null,
				  "domain": "{{.Domain}}",
				  "id": "{{ matchID64 }}",
				  "metadata": "*",
				  "publicName": "{{.PublicName}}",
				  "updatedAt": "{{ matchTimestamp }}",
				  "xpubId": "{{.XPubID}}"
				}`, map[string]any{
				"Address":    userPaymail,
				"PublicName": userPaymail.PublicName(),
				"Alias":      userPaymail.Alias(),
				"XPubID":     fixtures.Sender.XPubID(),
				"Domain":     fixtures.PaymailDomain,
			})
	})

	t.Run("try to search paymails as user", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForUser()

		// when:
		res, _ := client.R().Get("/api/v1/admin/paymails")

		// then:
		then.Response(res).IsUnauthorizedForUser()
	})
}

func TestPaymailLivecycle(t *testing.T) {
	// given:
	givenForAllTests := testabilities.Given(t)
	cleanup := givenForAllTests.StartedSPVWallet()
	defer cleanup()

	// and:
	newPaymail := fixtures.Paymail("newalias" + "@" + fixtures.PaymailDomain)

	var testState struct {
		newPaymailID          string
		paymailDetailsRawBody []byte
	}

	t.Run("add paymail for selected user as admin", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.R().
			SetBody(map[string]any{
				"key":        fixtures.Sender.XPub(),
				"address":    newPaymail,
				"publicName": newPaymail.PublicName(),
				"avatar":     "",
			}).
			Post("/api/v1/admin/paymails")

		// then:
		then.Response(res).
			HasStatus(http.StatusCreated).
			WithJSONMatching(`{
				"address": "{{.Address}}",
				"alias": "{{.Alias}}",
				"avatar": "",
				"createdAt": "{{ matchTimestamp }}",
				"deletedAt": null,
				"domain": "{{.Domain}}",
 				"id": "{{ matchID64 }}",
				"metadata": null,
				"publicName": "{{.PublicName}}",
				"updatedAt": "{{ matchTimestamp }}",
				"xpubId": "{{.XPubID}}"
			}`, map[string]any{
				"Address":    newPaymail,
				"PublicName": newPaymail.PublicName(),
				"Alias":      newPaymail.Alias(),
				"XPubID":     fixtures.Sender.XPubID(),
				"Domain":     fixtures.PaymailDomain,
			})

		// update:
		getter := then.Response(res).JSONValue()
		testState.newPaymailID = getter.GetString("id")
		testState.paymailDetailsRawBody = res.Body()
	})

	t.Run("get added paymail as admin", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.R().Get("/api/v1/admin/paymails/" + testState.newPaymailID)

		// then:
		then.Response(res).
			IsOK()

		// and:
		assert.Equal(t, testState.paymailDetailsRawBody, res.Body())
	})

	t.Run("remove paymail as admin", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.R().
			Delete("/api/v1/admin/paymails/" + testState.newPaymailID)

		// then:
		then.Response(res).IsOK()

		// verify paymail is deleted by trying to get it
		getRes, _ := client.R().Get("/api/v1/admin/paymails/" + testState.newPaymailID)
		then.Response(getRes).HasStatus(404)
	})

	t.Run("try to remove paymail as user", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForUser()

		// when:
		res, _ := client.R().
			Delete("/api/v1/admin/paymails/" + testState.newPaymailID)

		// then:
		then.Response(res).IsUnauthorizedForUser()
	})
}
