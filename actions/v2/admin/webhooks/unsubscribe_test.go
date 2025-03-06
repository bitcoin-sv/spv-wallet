package webhooks_test

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
)

func TestUnsubscribeWebhookHappyPath(t *testing.T) {
	t.Run("unsubscribe webhook", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(testabilities.Given(t), t)
		cleanup := given.StartedSPVWalletWithConfiguration(
			testengine.WithNotificationsEnabled(),
			testengine.WithV2())
		defer cleanup()
		client := given.HttpClient().ForAdmin()
		webhook := map[string]string{
			"url":         "http://localhost:8080",
			"tokenHeader": "Auth1",
			"tokenValue":  "123",
		}
		// when:
		res, _ := client.R().
			SetBody(webhook).
			Post(webhookURLSuffix)

		// then:
		then.Response(res).
			IsOK()

		// and:
		// when:
		res, _ = client.R().
			Get(webhookURLSuffix)

		// then:
		then.Response(res).
			IsOK().
			WithJSONf(`[
                {"url": "http://localhost:8080", "banned": false}
            ]`)

		// and:
		// when:
		res, _ = client.
			R().
			SetBody(
				map[string]string{
					"url": "http://localhost:8080",
				}).
			Delete(webhookURLSuffix)

		// then:
		then.Response(res).IsOK()

		// and:
		// when:
		res, _ = client.R().
			Get(webhookURLSuffix)

		// then:
		then.Response(res).
			IsOK().
			WithJSONf(`[
            ]`)
	})
}

func TestUnsubscribeWebhookErrorPath(t *testing.T) {
	t.Run("unsubscribe with invalid JSON returns bad request", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(testabilities.Given(t), t)
		cleanup := given.StartedSPVWalletWithConfiguration(
			testengine.WithNotificationsEnabled(),
			testengine.WithV2())
		defer cleanup()

		client := given.HttpClient().
			ForAdmin()

		// when:
		res, _ := client.
			R().
			SetBody("{invalid json}").
			Delete(webhookURLSuffix)

		// then:
		then.Response(res).
			HasStatus(400).
			WithJSONMatching(`{
				"code": "{{ .code }}",
				"message": "{{ .message }}"
			}`, map[string]any{
				"code":    spverrors.ErrCannotBindRequest.Code,
				"message": spverrors.ErrCannotBindRequest.Message,
			})
	})

	t.Run("unsubscribe with notification disabled returns error 404", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(testabilities.Given(t), t)
		cleanup := given.StartedSPVWalletWithConfiguration(
			testengine.WithV2())
		defer cleanup()
		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.
			R().
			SetBody(map[string]string{"url": "http://localhost:8080"}).
			Delete(webhookURLSuffix)

		// then:
		then.Response(res).
			HasStatus(404).
			WithJSONMatching(`{
				"code": "{{ .code }}",
				"message": "{{ .message }}"
			}`, map[string]any{
				"code":    spverrors.ErrNotificationsDisabled.Code,
				"message": spverrors.ErrNotificationsDisabled.Message,
			})
	})

	t.Run("unsubscribe non-existent webhook returns internal error", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(testabilities.Given(t), t)
		cleanup := given.StartedSPVWalletWithConfiguration(
			testengine.WithNotificationsEnabled(),
			testengine.WithV2())
		defer cleanup()

		client := given.HttpClient().ForAdmin()

		webhookURL := "http://nonexistent.com"

		// when:
		res, _ := client.
			R().
			SetBody(map[string]string{"url": webhookURL}).
			Delete(webhookURLSuffix)

		// then:
		then.Response(res).
			HasStatus(404).
			WithJSONMatching(`{
				"code": "{{ .code }}",
				"message": "{{ .message }}"
			}`, map[string]any{
				"code":    spverrors.ErrWebhookSubscriptionNotFound.Code,
				"message": spverrors.ErrWebhookSubscriptionNotFound.Message,
			})
	})

	t.Run("unsubscribe with missing URL returns bad request", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(testabilities.Given(t), t)
		cleanup := given.StartedSPVWalletWithConfiguration(
			testengine.WithNotificationsEnabled(),
			testengine.WithV2())
		defer cleanup()
		client := given.HttpClient().ForAdmin()

		// when:
		res, _ := client.
			R().
			SetBody(map[string]string{}).
			Delete(webhookURLSuffix)

		// then:
		then.Response(res).
			HasStatus(400).
			WithJSONMatching(`{
				"code": "{{ .code }}",
				"message": "{{ .message }}"
			}`, map[string]any{
				"code":    spverrors.ErrWebhookUrlRequired.Code,
				"message": spverrors.ErrWebhookUrlRequired.Message,
			})
	})
}
