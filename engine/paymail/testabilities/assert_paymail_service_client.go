package testabilities

import (
	"encoding/json"
	"testing"

	"github.com/bitcoin-sv/spv-wallet/engine/tester/jsonrequire"
	"github.com/bitcoin-sv/spv-wallet/engine/tester/paymailmock"
	"github.com/bitcoin-sv/spv-wallet/models/bsv"
	"github.com/stretchr/testify/require"
)

type PaymailExternalAssertions interface {
	ReceivedBeefTransaction(sender, beef, reference string)
	ReceivedP2PDestinationRequest(satoshis bsv.Satoshis) paymailmock.MockedP2PDestinationResponse
}

type PaymailCapabilityCallAssertions interface {
	WithRequestJSONMatching(expectedTemplateFormat string, params map[string]any)
}

func Then(t testing.TB, mockPaymail *paymailmock.PaymailClientMock) PaymailExternalAssertions {
	return &externalClientAssertions{
		t:           t,
		require:     require.New(t),
		mockPaymail: mockPaymail,
	}
}

type externalClientAssertions struct {
	t           testing.TB
	require     *require.Assertions
	mockPaymail *paymailmock.PaymailClientMock
}

func (e *externalClientAssertions) ReceivedBeefTransaction(sender, beef, reference string) {
	e.t.Helper()
	urlRegex := "beef"
	details := e.mockPaymail.GetCallByRegex(urlRegex)
	e.require.NotNil(details, "Expected call to %s", urlRegex)

	jsonrequire.Match(e.t, `
		{
			"beef": "{{ .beef }}",
			"decodedBeef": null,
			"hex": "",
			"metadata": {
				"sender": "{{ .sender }}"
			},
			"reference": "{{ .reference }}"
		}
	`, map[string]any{
		"beef":      beef,
		"sender":    sender,
		"reference": reference,
	}, string(details.RequestBody))
}

func (e *externalClientAssertions) ReceivedP2PDestinationRequest(requestedSatoshis bsv.Satoshis) paymailmock.MockedP2PDestinationResponse {
	e.t.Helper()
	urlRegex := "p2p-payment-destination"
	details := e.mockPaymail.GetCallByRegex(urlRegex)
	e.require.NotNil(details, "Expected call to %s", urlRegex)

	jsonrequire.Match(e.t, `{
		"satoshis": {{ .satoshis }}
	}`, map[string]any{"satoshis": requestedSatoshis}, string(details.RequestBody))

	var response paymailmock.MockedP2PDestinationResponse
	err := json.Unmarshal(details.ResponseBody, &response)
	e.require.NoError(err)

	return response
}
