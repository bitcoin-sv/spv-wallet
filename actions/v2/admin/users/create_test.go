package users_test

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	"github.com/bitcoin-sv/spv-wallet/actions/testabilities/apierror"
	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/fixtures"
)

func TestCreateUserWithoutPaymail(t *testing.T) {
	// given:
	givenForAllTests := testabilities.Given(t)
	cleanup := givenForAllTests.StartedSPVWalletWithConfiguration(
		testengine.WithV2(),
	)
	defer cleanup()

	// and:
	userCandidate := fixtures.User{
		PrivKey: "xprv9s21ZrQH143K31pvNoYNcRZjtdJXnNVEc5NmBbgJmEg27YWbZVL7jTLQhPELqAR7tcJTnF9AJLwVN5w3ABZvrfeDLm4vnBDw76bkx8a2NxK",
	}
	publicKey := userCandidate.PublicKey().ToDERHex()

	var testState struct {
		userID string
	}

	t.Run("Create a user as admin", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.R().
			SetBody(map[string]any{
				"publicKey": publicKey,
			}).
			Post("/api/v2/admin/users")

		// then:
		then.Response(res).
			HasStatus(201).
			WithJSONMatching(`{
				"id": "{{ matchAddress }}",
				"createdAt": "{{ matchTimestamp }}",
				"updatedAt": "{{ matchTimestamp }}",
				"publicKey": "{{ .publicKey }}",
				"paymails": []
			}`, map[string]any{
				"publicKey": publicKey,
			})

		// update:
		getter := then.Response(res).JSONValue()
		testState.userID = getter.GetString("id")
	})

	t.Run("Get new user by id as admin", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.R().
			SetPathParam("id", testState.userID).
			Get("/api/v2/admin/users/{id}")

		then.Response(res).
			IsOK().
			WithJSONMatching(`{
				"id": "{{ matchAddress }}",
				"createdAt": "{{ matchTimestamp }}",
				"updatedAt": "{{ matchTimestamp }}",
				"publicKey": "{{ .publicKey }}",
				"paymails": []
			}`, map[string]any{
				"publicKey": publicKey,
			})
	})

	t.Run("Get user info as a new user", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForGivenUser(userCandidate)

		// when:
		res, _ := client.R().Get("/api/v2/users/current")

		// then:
		then.Response(res).
			IsOK().
			WithJSONMatching(`{
				"currentBalance": 0
			}`, nil)
	})
}

func TestCreateUserWithPaymail(t *testing.T) {
	// given:
	givenForAllTests := testabilities.Given(t)
	cleanup := givenForAllTests.StartedSPVWalletWithConfiguration(
		testengine.WithV2(),
	)
	defer cleanup()

	// and:
	userCandidate := fixtures.User{
		PrivKey: "xprv9s21ZrQH143K31pvNoYNcRZjtdJXnNVEc5NmBbgJmEg27YWbZVL7jTLQhPELqAR7tcJTnF9AJLwVN5w3ABZvrfeDLm4vnBDw76bkx8a2NxK",
		Paymails: []fixtures.Paymail{
			"test_user@" + fixtures.PaymailDomain,
		},
	}
	publicKey := userCandidate.PublicKey().ToDERHex()
	publicName := "Test User"

	var testState struct {
		userID string
	}

	t.Run("Create a user as admin", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.R().
			SetBody(map[string]any{
				"publicKey": publicKey,
				"paymail": map[string]any{
					"address":    userCandidate.DefaultPaymail(),
					"publicName": publicName,
					"avatar":     "",
				},
			}).
			Post("/api/v2/admin/users")

		// then:
		then.Response(res).
			HasStatus(201).
			WithJSONMatching(`{
				"id": "{{ matchAddress }}",
				"createdAt": "{{ matchTimestamp }}",
				"updatedAt": "{{ matchTimestamp }}",
				"publicKey": "{{ .publicKey }}",
				"paymails": [
					{
						"alias": "{{ .alias }}",
						"avatar": "",
						"domain": "example.com",
						"id": "{{ matchNumber }}",
						"paymail": "{{ .paymail }}",
						"publicName": "{{ .publicName }}"
					}
				]
			}`, map[string]any{
				"publicKey":  publicKey,
				"paymail":    userCandidate.DefaultPaymail(),
				"publicName": publicName,
				"alias":      userCandidate.DefaultPaymail().Alias(),
			})

		// update:
		getter := then.Response(res).JSONValue()
		testState.userID = getter.GetString("id")
	})

	t.Run("Get new user by id as admin", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.R().
			SetPathParam("id", testState.userID).
			Get("/api/v2/admin/users/{id}")

		then.Response(res).
			IsOK().
			WithJSONMatching(`{
				"id": "{{ matchAddress }}",
				"createdAt": "{{ matchTimestamp }}",
				"updatedAt": "{{ matchTimestamp }}",
				"publicKey": "{{ .publicKey }}",
				"paymails": [
					{
						"alias": "{{ .alias }}",
						"avatar": "",
						"domain": "example.com",
						"id": "{{ matchNumber }}",
						"paymail": "{{ .paymail }}",
						"publicName": "{{ .publicName }}"
					}
				]
			}`, map[string]any{
				"publicKey":  publicKey,
				"paymail":    userCandidate.DefaultPaymail(),
				"publicName": publicName,
				"alias":      userCandidate.DefaultPaymail().Alias(),
			})
	})
}

func TestCreateUserWithAliasAndDomain(t *testing.T) {
	// given:
	givenForAllTests := testabilities.Given(t)
	cleanup := givenForAllTests.StartedSPVWalletWithConfiguration(
		testengine.WithV2(),
	)
	defer cleanup()

	// and:
	userCandidate := fixtures.User{
		PrivKey: "xprv9s21ZrQH143K31pvNoYNcRZjtdJXnNVEc5NmBbgJmEg27YWbZVL7jTLQhPELqAR7tcJTnF9AJLwVN5w3ABZvrfeDLm4vnBDw76bkx8a2NxK",
		Paymails: []fixtures.Paymail{
			"test_user@" + fixtures.PaymailDomain,
		},
	}
	publicKey := userCandidate.PublicKey().ToDERHex()
	publicName := "Test User"

	var testState struct {
		userID string
	}

	t.Run("Create a user as admin", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.R().
			SetBody(map[string]any{
				"publicKey": publicKey,
				"paymail": map[string]any{
					"alias":      userCandidate.DefaultPaymail().Alias(),
					"domain":     fixtures.PaymailDomain,
					"publicName": publicName,
					"avatar":     "",
				},
			}).
			Post("/api/v2/admin/users")

		// then:
		then.Response(res).
			HasStatus(201).
			WithJSONMatching(`{
				"id": "{{ matchAddress }}",
				"createdAt": "{{ matchTimestamp }}",
				"updatedAt": "{{ matchTimestamp }}",
				"publicKey": "{{ .publicKey }}",
				"paymails": [
					{
						"alias": "{{ .alias }}",
						"avatar": "",
						"domain": "example.com",
						"id": "{{ matchNumber }}",
						"paymail": "{{ .paymail }}",
						"publicName": "{{ .publicName }}"
					}
				]
			}`, map[string]any{
				"publicKey":  publicKey,
				"paymail":    userCandidate.DefaultPaymail(),
				"publicName": publicName,
				"alias":      userCandidate.DefaultPaymail().Alias(),
			})

		// update:
		getter := then.Response(res).JSONValue()
		testState.userID = getter.GetString("id")
	})

	t.Run("Get new user by id as admin", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(givenForAllTests, t)
		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.R().
			SetPathParam("id", testState.userID).
			Get("/api/v2/admin/users/{id}")

		then.Response(res).
			IsOK().
			WithJSONMatching(`{
				"id": "{{ matchAddress }}",
				"createdAt": "{{ matchTimestamp }}",
				"updatedAt": "{{ matchTimestamp }}",
				"publicKey": "{{ .publicKey }}",
				"paymails": [
					{
						"alias": "{{ .alias }}",
						"avatar": "",
						"domain": "example.com",
						"id": "{{ matchNumber }}",
						"paymail": "{{ .paymail }}",
						"publicName": "{{ .publicName }}"
					}
				]
			}`, map[string]any{
				"publicKey":  publicKey,
				"paymail":    userCandidate.DefaultPaymail(),
				"publicName": publicName,
				"alias":      userCandidate.DefaultPaymail().Alias(),
			})
	})
}

func TestAddUserWithWrongPaymailDomain(t *testing.T) {
	// given:
	given, then := testabilities.New(t)
	cleanup := given.StartedSPVWalletWithConfiguration(
		testengine.WithV2(),
	)
	defer cleanup()

	// and:
	userCandidate := fixtures.User{
		PrivKey: "xprv9s21ZrQH143K31pvNoYNcRZjtdJXnNVEc5NmBbgJmEg27YWbZVL7jTLQhPELqAR7tcJTnF9AJLwVN5w3ABZvrfeDLm4vnBDw76bkx8a2NxK",
	}
	publicKey := userCandidate.PublicKey().ToDERHex()

	// and:
	alias := "user"
	unsupportedDomain := "unsupported.com"
	unsupportedPaymail := alias + "@" + unsupportedDomain
	publicName := "User"

	// and:
	client := given.HttpClient().ForAdmin()

	t.Run("Try to add using whole paymail as address", func(t *testing.T) {
		// when:
		res, _ := client.R().
			SetBody(map[string]any{
				"publicKey": publicKey,
				"paymail": map[string]any{
					"address":    unsupportedPaymail,
					"publicName": publicName,
					"avatar":     "",
				},
			}).
			Post("/api/v2/admin/users")

		// then:
		then.Response(res).
			HasStatus(400).
			WithJSONf(apierror.ExpectedJSON("error-invalid-domain", "invalid domain"))
	})

	t.Run("Try to add using alias and domain as address", func(t *testing.T) {
		// when:
		res, _ := client.R().
			SetBody(map[string]any{
				"publicKey": publicKey,
				"paymail": map[string]any{
					"alias":      alias,
					"domain":     unsupportedDomain,
					"publicName": publicName,
					"avatar":     "",
				},
			}).
			Post("/api/v2/admin/users")

		// then:
		then.Response(res).
			HasStatus(400).
			WithJSONf(apierror.ExpectedJSON("error-invalid-domain", "invalid domain"))
	})

}
