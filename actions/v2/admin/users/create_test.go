package users_test

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
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

func TestCreateUserWithBadURLAvatar(t *testing.T) {
	// given:
	given, then := testabilities.New(t)
	cleanup := given.StartedSPVWalletWithConfiguration(
		testengine.WithV2(),
	)
	defer cleanup()

	// and:
	client := given.HttpClient().ForAdmin()

	// and:
	userCandidate := fixtures.User{
		PrivKey: "xprv9s21ZrQH143K4Tf1hf7ouMiagMH4JKvE6E2SY8Su55Y6aFi9AfQibzx7i79g1NJkLQbRY4FjKgvpddtYXoD7dA2KbGjHdHcxXVqtd687KrK",
		Paymails: []fixtures.Paymail{
			"second@" + fixtures.PaymailDomain,
		},
	}
	publicKey := userCandidate.PublicKey().ToDERHex()
	avatarURL := "User/Path/To/Avatar"

	// when:
	res, _ := client.R().
		SetBody(map[string]any{
			"publicKey": publicKey,
			"paymail": map[string]any{
				"address":   userCandidate.DefaultPaymail(),
				"avatarURL": avatarURL,
			},
		}).
		Post("/api/v2/admin/users")

	// then:
	then.Response(res).
		WithProblemDetails(422, "invalid_avatar_url", "Invalid avatar URL")

}

func TestCreateUserWithoutPublicName(t *testing.T) {
	// given:
	given, then := testabilities.New(t)
	cleanup := given.StartedSPVWalletWithConfiguration(testengine.WithV2())
	defer cleanup()

	// and:
	client := given.HttpClient().ForAdmin()

	// and:
	userCandidate := fixtures.User{
		PrivKey: "xprv9s21ZrQH143K4Tf1hf7ouMiagMH4JKvE6E2SY8Su55Y6aFi9AfQibzx7i79g1NJkLQbRY4FjKgvpddtYXoD7dA2KbGjHdHcxXVqtd687KrK",
		Paymails: []fixtures.Paymail{
			"second@" + fixtures.PaymailDomain,
		},
	}
	publicKey := userCandidate.PublicKey().ToDERHex()
	avatarURL := "https://address-to-avatar.com"

	// when:
	res, _ := client.R().
		SetBody(map[string]any{
			"publicKey": publicKey,
			"paymail": map[string]any{
				"address":   userCandidate.DefaultPaymail(),
				"avatarURL": avatarURL,
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
						"avatar": "{{ .avatar }}",
						"domain": "example.com",
						"id": "{{ matchNumber }}",
						"paymail": "{{ .paymail }}",
						"publicName": "{{ .alias }}"
					}
				]
			}`, map[string]any{
			"publicKey": publicKey,
			"paymail":   userCandidate.DefaultPaymail(),
			"alias":     userCandidate.DefaultPaymail().Alias(),
			"avatar":    avatarURL,
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
					"avatarURL":  "",
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
					"avatarURL":  "",
				},
			}).
			Post("/api/v2/admin/users")

		// then:
		then.Response(res).
			WithProblemDetails(400, "unsupported_domain", "Unsupported domain")
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
					"avatarURL":  "",
				},
			}).
			Post("/api/v2/admin/users")

		// then:
		then.Response(res).
			WithProblemDetails(400, "unsupported_domain", "Unsupported domain")
	})
}

func TestTryToAddWithWrongPubKey(t *testing.T) {
	// given:
	given, then := testabilities.New(t)
	cleanup := given.StartedSPVWalletWithConfiguration(
		testengine.WithV2(),
	)
	defer cleanup()

	// and:
	client := given.HttpClient().ForAdmin()

	// when:
	res, _ := client.R().
		SetBody(map[string]any{
			"publicKey": "wrong",
		}).
		Post("/api/v2/admin/users")

	// then:
	then.Response(res).
		WithProblemDetails(400, "invalid_public_key", "Invalid public key")
}
