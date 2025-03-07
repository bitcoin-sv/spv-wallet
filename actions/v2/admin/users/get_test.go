package users_test

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
)

func TestGetUser(t *testing.T) {
	// given:
	given, then := testabilities.New(t)
	cleanup := given.StartedSPVWalletWithConfiguration(testengine.WithV2())
	defer cleanup()

	// and:
	client := given.HttpClient().ForAdmin()

	// when:
	resp, _ := client.
		R().
		SetPathParam("id", fixtures.Sender.ID()).
		Get("/api/v2/admin/users/{id}")

	// then:
	then.Response(resp).
		HasStatus(200).
		WithJSONMatching(`{
			"id": "{{ .id }}",
			"createdAt": "{{ matchTimestamp }}",
			"updatedAt": "{{ matchTimestamp }}",
			"publicKey": "{{ .publicKey }}",
			"paymails": [
			    {
			      "alias": "{{ .alias }}",
			      "avatar": "",
			      "domain": "{{ .domain }}",
			      "id": 3,
			      "paymail": "{{ .paymail }}",
			      "publicName": "{{ .publicName }}"
			    }
			]
		}`, map[string]any{
			"id":         fixtures.Sender.ID(),
			"publicKey":  fixtures.Sender.PublicKey().ToDERHex(),
			"paymail":    fixtures.Sender.DefaultPaymail(),
			"publicName": fixtures.Sender.DefaultPaymail().PublicName(),
			"alias":      fixtures.Sender.DefaultPaymail().Alias(),
			"domain":     fixtures.Sender.DefaultPaymail().Domain(),
		})
}

func TestTryGetNonExistingUser(t *testing.T) {
	// given:
	given, then := testabilities.New(t)
	cleanup := given.StartedSPVWalletWithConfiguration(testengine.WithV2())
	defer cleanup()

	// and:
	client := given.HttpClient().ForAdmin()

	// when:
	resp, _ := client.
		R().
		SetPathParam("id", "non-existing-id").
		Get("/api/v2/admin/users/{id}")

	// then:
	then.Response(resp).
		WithProblemDetails(404, "not_found")
}
