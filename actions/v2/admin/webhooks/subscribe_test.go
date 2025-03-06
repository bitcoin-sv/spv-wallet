package webhooks_test

import (
	"testing"

	"github.com/bitcoin-sv/spv-wallet/actions/testabilities"
	"github.com/bitcoin-sv/spv-wallet/engine/spverrors"
	testengine "github.com/bitcoin-sv/spv-wallet/engine/testabilities"
)

const webhookURLSuffix = "/api/v2/admin/webhooks"

func TestSubscribeWebhooksHappyPath(t *testing.T) {
	t.Run("add webhook", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(testabilities.Given(t), t)
		cleanup := given.StartedSPVWalletWithConfiguration(
			testengine.WithNotificationsEnabled(),
			testengine.WithV2())
		defer cleanup()

		client := given.HttpClient().
			ForAdmin()

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
	})

	t.Run("add same webhook 2 times", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(testabilities.Given(t), t)
		cleanup := given.StartedSPVWalletWithConfiguration(
			testengine.WithNotificationsEnabled(),
			testengine.WithV2())
		defer cleanup()

		client := given.HttpClient().
			ForAdmin()

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
	})

	t.Run("subscribe and retrieve multiple webhooks", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(testabilities.Given(t), t)
		cleanup := given.StartedSPVWalletWithConfiguration(
			testengine.WithNotificationsEnabled(),
			testengine.WithV2())
		defer cleanup()

		client := given.HttpClient().
			ForAdmin()

		webhook1 := map[string]string{
			"url":         "http://localhost:8080",
			"tokenHeader": "Auth1",
			"tokenValue":  "123",
		}

		webhook2 := map[string]string{
			"url":         "http://localhost:8081",
			"tokenHeader": "Auth2",
			"tokenValue":  "456",
		}

		// when:
		client.R().
			SetBody(webhook1).
			Post(webhookURLSuffix)

		client.R().
			SetBody(webhook2).
			Post(webhookURLSuffix)

		res, _ := client.R().
			Get(webhookURLSuffix)

		// then:
		then.Response(res).
			IsOK().
			WithJSONf(`[
                {"url": "http://localhost:8080", "banned": false},
                {"url": "http://localhost:8081", "banned": false}
            ]`)
	})

}

func TestSubscribeWebhooksErrorPath(t *testing.T) {
	t.Run("add webhook when notifications disabled", func(t *testing.T) {
		// given:
		given, then := testabilities.NewOf(testabilities.Given(t), t)
		cleanup := given.StartedSPVWalletWithConfiguration(
			testengine.WithV2())
		defer cleanup()

		client := given.HttpClient().
			ForAdmin()

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
			HasStatus(404).
			WithJSONMatching(`{
				"code": "{{ .code }}",
				"message": "{{ .message }}"
			}`, map[string]any{
				"code":    spverrors.ErrNotificationsDisabled.Code,
				"message": spverrors.ErrNotificationsDisabled.Message,
			})
	})

	t.Run("subscribe with invalid JSON returns bad request", func(t *testing.T) {
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
			Post(webhookURLSuffix)

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

	t.Run("subscribe with missing token value returns bad request", func(t *testing.T) {
		// given:
		given, then := testabilities.New(t)
		cleanup := given.StartedSPVWalletWithConfiguration(
			testengine.WithNotificationsEnabled(),
			testengine.WithV2())
		defer cleanup()

		client := given.HttpClient().
			ForAdmin()

		webhook := map[string]string{
			"tokenHeader": "Authorization",
			"url":         "http://localhost:8080",
		}

		// when:
		res, _ := client.
			R().
			SetBody(webhook).
			Post(webhookURLSuffix)

		// then:
		then.Response(res).
			HasStatus(400).
			WithJSONMatching(`{
				"code": "{{ .code }}",
				"message": "{{ .message }}"
			}`, map[string]any{
				"code":    spverrors.ErrWebhookTokenValueRequired.Code,
				"message": spverrors.ErrWebhookTokenValueRequired.Message,
			})

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

	t.Run("subscribe with missing token header returns bad request", func(t *testing.T) {
		// given:
		given, then := testabilities.New(t)
		cleanup := given.StartedSPVWalletWithConfiguration(
			testengine.WithNotificationsEnabled(),
			testengine.WithV2())
		defer cleanup()

		client := given.HttpClient().
			ForAdmin()

		webhook := map[string]string{
			"url":        "http://localhost:8080",
			"tokenValue": "123",
		}

		// when:
		res, _ := client.
			R().
			SetBody(webhook).
			Post(webhookURLSuffix)

		// then:
		then.Response(res).
			HasStatus(400).
			WithJSONMatching(`{
				"code": "{{ .code }}",
				"message": "{{ .message }}"
			}`, map[string]any{
				"code":    spverrors.ErrWebhookTokenHeaderRequired.Code,
				"message": spverrors.ErrWebhookTokenHeaderRequired.Message,
			})

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

	t.Run("subscribe with missing URL returns bad request", func(t *testing.T) {
		// given:
		given, then := testabilities.New(t)
		cleanup := given.StartedSPVWalletWithConfiguration(
			testengine.WithNotificationsEnabled(),
			testengine.WithV2())
		defer cleanup()

		client := given.HttpClient().
			ForAdmin()

		webhook := map[string]string{
			"tokenHeader": "Authorization",
			"tokenValue":  "123",
		}

		// when:
		res, _ := client.
			R().
			SetBody(webhook).
			Post(webhookURLSuffix)

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
